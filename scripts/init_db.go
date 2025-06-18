package main

import (
    "log"

    "github.com/start-finish/startfront-app/internal/models"
    "github.com/start-finish/startfront-app/pkg/database"
)

func main() {
    db := database.Init()

    err := db.AutoMigrate(
        &models.User{}, // Add more models as needed
    )
    if err != nil {
        log.Fatal("❌ AutoMigrate failed:", err)
    }

    log.Println("✅ Database migration completed.")
}
