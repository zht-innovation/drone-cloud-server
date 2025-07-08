package main

import (
	"zhtcloud/gateway"
	"zhtcloud/models"
	"zhtcloud/utils/logger"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file: %v", err)
		return
	}

	models.DatabaseSetup()
	gateway.ServerSetup()

	defer models.Close()
}
