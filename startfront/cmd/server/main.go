package main

import (
	"fmt"
	"log"

	"startfront/config"
	"startfront/docs"
	"startfront/models"
	"startfront/router"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// migrate base models
	if err := db.Migrator().AutoMigrate(
		&models.Users{},
		&models.Clients{},
		&models.Projects{},
		&models.Client_users{},
		&models.Screens{},
		&models.Navigation_menus{},
		&models.Navigation_items{},
		&models.Widgets{},
		&models.Widget_instances{},
		&models.Screen_widgets{},
		&models.Widget_presets{},
		&models.Widget_preset_items{},
		&models.Themes{},
		&models.Settings{},
		&models.Analytics_events{},
		&models.Roles{},
		&models.User_roles{},
		&models.Permissions{},
		&models.Role_permissions{},
		&models.Admin_dashboard{},
	); err != nil {
		log.Fatal(err)
	}

	// Seed database with default roles and users
	if err := seedDatabase(db); err != nil {
		log.Fatal("❌ Error seeding database:", err)
	}

	r := gin.Default()
	router.SetupRoutes(r, db)
	docs.SetupSwagger(r)

	fmt.Println("📚 Swagger UI: http://localhost:8080/swagger/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func seedDatabase(db *gorm.DB) error {
	// Import models if needed, they are already imported as models.Roles etc.
	var rolesCount int64
	db.Model(&models.Roles{}).Count(&rolesCount)
	if rolesCount == 0 {
		defaultRoles := []models.Roles{
			{RoleName: "Super Admin", Description: "Full access to all platform features", UserCount: 2, Permissions: `["read","write","delete","admin"]`},
			{RoleName: "Admin", Description: "Administrative access with some restrictions", UserCount: 1, Permissions: `["read","write","delete"]`},
			{RoleName: "Editor", Description: "Can create and edit content", UserCount: 1, Permissions: `["read","write"]`},
			{RoleName: "Viewer", Description: "Read-only access", UserCount: 1, Permissions: `["read"]`},
		}
		for _, r := range defaultRoles {
			if err := db.Create(&r).Error; err != nil {
				return err
			}
		}
		fmt.Println("🌱 Seeded default roles successfully")
	}

	var usersCount int64
	db.Model(&models.Users{}).Count(&usersCount)
	if usersCount == 0 {
		hashedAlex, _ := bcrypt.GenerateFromPassword([]byte("alex123"), bcrypt.DefaultCost)
		hashedSarah, _ := bcrypt.GenerateFromPassword([]byte("sarah123"), bcrypt.DefaultCost)
		hashedMike, _ := bcrypt.GenerateFromPassword([]byte("mike123"), bcrypt.DefaultCost)
		hashedLisa, _ := bcrypt.GenerateFromPassword([]byte("lisa123"), bcrypt.DefaultCost)
		hashedJames, _ := bcrypt.GenerateFromPassword([]byte("james123"), bcrypt.DefaultCost)

		defaultUsers := []models.Users{
			{Username: "Alex Johnson", Email: "alex@startfront.io", Password: string(hashedAlex), Role: "Super Admin", Description: "Full access to all platform features and settings.", Status: "Active"},
			{Username: "Sarah Chen", Email: "sarah@startfront.io", Password: string(hashedSarah), Role: "Admin", Description: "Administrative access with some system restrictions.", Status: "Active"},
			{Username: "Mike Torres", Email: "mike@startfront.io", Password: string(hashedMike), Role: "Viewer", Description: "Read-only access to dashboard and reporting tools.", Status: "Inactive"},
			{Username: "Lisa Wang", Email: "lisa@startfront.io", Password: string(hashedLisa), Role: "Editor", Description: "Can create, edit and manage platform content.", Status: "Active"},
			{Username: "James Park", Email: "james@startfront.io", Password: string(hashedJames), Role: "Viewer", Description: "Read-only access to dashboard and reporting tools.", Status: "Active"},
		}
		for _, u := range defaultUsers {
			if err := db.Create(&u).Error; err != nil {
				return err
			}
		}
		fmt.Println("🌱 Seeded default users successfully")
	}

	var settingsCount int64
	db.Model(&models.Settings{}).Count(&settingsCount)
	if settingsCount == 0 {
		defaultSettings := []models.Settings{
			{Key: "Platform Name", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"StartFront"}`))},
			{Key: "Default Language", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"English"}`))},
			{Key: "Timezone", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"UTC-8"}`))},
			{Key: "Connection Pool Size", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"20"}`))},
			{Key: "Query Timeout", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"30s"}`))},
			{Key: "Backup Frequency", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"Daily"}`))},
			{Key: "SMTP Server", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"smtp.startfront.com"}`))},
			{Key: "From Address", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"noreply@startfront.com"}`))},
			{Key: "Daily Limit", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"10000"}`))},
			{Key: "Session Timeout", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"24 hours"}`))},
			{Key: "Password Policy", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"Strong"}`))},
			{Key: "2FA Required", Scope: "global", Value: datatypes.JSON([]byte(`{"value":"false"}`))},
		}
		for _, s := range defaultSettings {
			if err := db.Create(&s).Error; err != nil {
				return err
			}
		}
		fmt.Println("🌱 Seeded default settings successfully")
	}

	return nil
}
