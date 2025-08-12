package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"startfront/config"
	"startfront/models"
	"startfront/router"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// migrate base models
	if err := db.Migrator().AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	router.SetupRoutes(r, db)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
