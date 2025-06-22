// File: cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/start-finish/startfront-app/internal/boot"
	"github.com/start-finish/startfront-app/internal/engine"
	"github.com/start-finish/startfront-app/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using system environment variables")
	}

	// Connect to database
	db := database.Init()

	// Run auto migration for core static models
	boot.AutoMigrate(db)

	// Initialize Gin router
	router := gin.Default()

	// // Register core application routes
	// auth.RegisterRoutes(router, db)
	// user.RegisterRoutes(router, db)

	// Load dynamic schemas
	schemas, err := engine.LoadSchemas("./schemas/")
	if err != nil {
		log.Fatalf("❌ Failed to load schemas: %v", err)
	}

	// Auto-migrate and register routes for each dynamic schema
	for _, schema := range schemas {
		if err := engine.AutoMigrateSchema(db, schema); err != nil {
			log.Fatalf("❌ Auto-migrate failed for schema %s: %v", schema.Model, err)
		}

		engine.RegisterRoutes(router, db, schema)
	}

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Get IP and Port
	ip := os.Getenv("SERVER_IP")
	if ip == "" {
		ip = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Print all registered routes
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
