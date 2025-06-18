package database

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env file not found, using system environment variables.")
    }

    // Use environment variables or fallback to config.yaml
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")

    if dbHost == "" {
        dbHost = "localhost"
    }
    if dbPort == "" {
        dbPort = "5432"
    }
    if dbUser == "" {
        dbUser = "startadmin"
    }
    if dbPass == "" {
        dbPass = "startpassword"
    }
    if dbName == "" {
        dbName = "startfrontdb"
    }

    // Format the DSN for PostgreSQL
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

    // Connect to the database
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("❌ Failed to connect to database:", err)
    }

    DB = db
    fmt.Println("✅ Connected to database")
    return db
}
