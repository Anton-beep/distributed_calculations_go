package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"storage/internal/Db"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		zap.S().Fatal(err)
	}

	_, err = Db.New()
	if err != nil {
		zap.S().Fatal(err)
	}
}
