package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Add more modules here
	"github.com/start-finish/startfront-app/internal/users"
	"github.com/start-finish/startfront-app/pkg"
)

func main() {
	// 1️⃣ Load .env first
	godotenv.Load()

	// 2️⃣ Connect DB
	db := pkg.Connect()

	// 3️⃣ Register Modules
	pkg.Modules = []pkg.Module{
		// &projects.ProjectModule{},
		&users.UserModule{},
	}

	// 4️⃣ Auto Migrate All Modules
	for _, m := range pkg.Modules {
		m.AutoMigrate(db)
	}

	// 5️⃣ Init router
	router := gin.Default()

	// 6️⃣ Register Routes for All Modules
	pkg.RegisterModules(router, db)

	// 8️⃣ Print All Routes
	printRoutes(router)

	// 9️⃣ Run server
	router.Run(":8080")
}

func printRoutes(r *gin.Engine) {
	for _, route := range r.Routes() {
		fmt.Printf("✅ %-6s %s\n", route.Method, route.Path)
	}
}
