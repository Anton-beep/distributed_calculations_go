package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	// gin.SetMode(gin.ReleaseMode)

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

	// api
	err = api.New(d).Start().Run(":8080")
	if err != nil {
		zap.S().Fatal(err)
	}
}
