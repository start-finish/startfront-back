package boot

import (
	"log"

	"github.com/start-finish/startfront-app/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	// List of models (tables) to migrate
	modelsToMigrate := []interface{}{
		&models.User{},
		// Add other models here as needed
		// &models.Session{},
		// &models.Application{},
		// &models.ApplicationCollaborator{},
	}

	// Iterate over each model and perform migration dynamically
	for _, model := range modelsToMigrate {
		// We use the GORM migrator's AutoMigrate method to check if the table exists
		if !db.Migrator().HasTable(model) {
			// If table doesn't exist, run the migration
			log.Printf("✅ Table for '%T' does not exist, running migration...", model)
			if err := db.AutoMigrate(model); err != nil {
				log.Fatalf("❌ Auto migration failed for '%T': %v", model, err)
			}
			log.Printf("✅ Table for '%T' migration completed.", model)
		} else {
			log.Printf("✅ Table for '%T' already exists, skipping migration.", model)
		}
	}
}