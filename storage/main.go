package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	_ "storage/docs"
	"storage/internal/Api"
	"storage/internal/Db"
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
// @BasePath	/api/v1
func main() {
	InitLogger(true)
	//gin.SetMode(gin.ReleaseMode)

	// .env
	err := godotenv.Load()
	if err != nil {
		zap.S().Fatal(err)
	}

	// db
	d, err := Db.New()
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	// api
	r := Api.New(d)
	err = r.Start().Run(":8080")
	if err != nil {
		zap.S().Fatal(err)
	}
}
