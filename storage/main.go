package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"storage/internal/availableservers"
	"storage/internal/expressionstorage"
	"storage/internal/frontend"
	"storage/internal/gRPCServer"
	"strconv"
	"sync"
	"time"

	// postgresql driver.
	_ "storage/docs"
	"storage/internal/api"
	"storage/internal/db"
)

func InitLogger(debug bool) {
	cfg := zap.NewDevelopmentConfig()
	if debug {
		cfg.Level.SetLevel(zap.DebugLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	zap.S().Info("Start")
}

//	@title			Swagger Storage API
//	@version		1.0
//	@description	This is a server for the storage of expressions and their results

// @host		localhost:8080
// @BasePath	/api/v1.
func main() {
	InitLogger(true)
	gin.SetMode(gin.ReleaseMode)

	// .env
	err := godotenv.Load()
	if err != nil {
		zap.S().Warn(err)
	}

	// db
	d, err := db.New()
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	// workers storage
	workerStorage := sync.Map{}

	// expression storage
	num, err := strconv.Atoi(os.Getenv("CHECK_SERVER_DURATION"))
	expStorage := expressionstorage.New(d, time.Duration(num)*time.Second, &workerStorage)

	// servers storage
	servers := availableservers.New(expStorage)

	// execution time configs
	execTimeConfig := &api.ExecTimeConfig{}

	// frontend build
	go frontend.ServeFrontend()

	server := gRPCServer.New(expStorage, servers, execTimeConfig, &workerStorage) // create a new gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051") // specify the port you want your gRPC server to run on
		if err != nil {
			zap.S().Fatal(err)
		}
		grpcServer := grpc.NewServer()
		gRPCServer.RegisterExpressionsServiceServer(grpcServer, server)
		if err := grpcServer.Serve(lis); err != nil {
			zap.S().Fatal(err)
		}
	}()

	// api
	err = api.New(d, expStorage, &workerStorage, servers, execTimeConfig).Start().Run(":8080")
	if err != nil {
		zap.S().Fatal(err)
	}
}
