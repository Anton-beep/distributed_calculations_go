package main

import (
	_ "calculationServer/docs"
	"calculationServer/internal/storageclient"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

func main() {
	InitLogger(true)

	// .env
	err := godotenv.Load()
	if err != nil {
		zap.S().Warn(err)
	}

	// start client work
	c, err := storageclient.New()
	if err != nil {
		zap.S().Fatal(err)
	}
	err = c.SetupgRPCServer()
	if err != nil {
		zap.S().Fatal(err)
	}

	c.Run()
}
