package main

import (
	"startfront-backend/internal/delivery"
	"startfront-backend/internal/repository"
	"startfront-backend/pkg/config"
)

func main() {
	config.InitDB()
	repository.InitDB(config.DBConfigData)
	delivery.StartServer(config.DB)
}
