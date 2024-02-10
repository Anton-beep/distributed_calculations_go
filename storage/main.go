package main

import (
	"github.com/joho/godotenv"
	"log"
	"storage/internal/Db"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.New()
	if err != nil {
		return
	}
}
