// File: cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/start-finish/startfront-app/internal/auth"
	"github.com/start-finish/startfront-app/internal/boot"
	"github.com/start-finish/startfront-app/internal/user"
	"github.com/start-finish/startfront-app/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using system environment variables")
	}

	// Connect to database
	db := database.Init()

	// Drop all tables
	// boot.DropAllTables(db)

	// Run auto migration from centralized boot package
	boot.AutoMigrate(db)

	// Initialize Gin router
	router := gin.Default()

	// Register application routes
	auth.RegisterRoutes(router, db)
	user.RegisterRoutes(router, db)

	// List all routes for debugging
	// printRoutes(router)

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Get the IP and port
	ip := os.Getenv("SERVER_IP")
	if ip == "" {
		ip = "localhost" // Default to localhost if not set
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}

	// Print all registered routes with full URL
	printRoutes(router, ip, port)

	// Start server
	fmt.Printf("🚀 Server running on http://%s:%s\n", ip, port)
	router.Run(":" + port)
}

// Function to print all registered routes with full URL
func printRoutes(router *gin.Engine, ip, port string) {
	routes := router.Routes()
	fmt.Println("Registered Routes:")
	for _, route := range routes {
		fullURL := fmt.Sprintf("http://%s:%s%s", ip, port, route.Path)
		fmt.Printf("Method: %-6s | Full Route: %-50s\n", route.Method, fullURL)
	}
}
