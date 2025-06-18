package boot

import (
	"log"
	"gorm.io/gorm"
)

func DropAllTables(db *gorm.DB) {
	// List of all tables to drop
	tables := []string{
		"users",
		// "applications",
		// "application_collaborators",
		// "app_connections",
		// "screens",
		// "widgets",
		// "widget_bindings",
		// "widget_presets",
		// "widget_props",
		// "application_versions",
		// "navigation_menus",
		// "assets",
		// "feedback",
		// "audit_logs",
		// "analytics_logs",
		// "export_reports",
		// "themes",
		// "settings",
		// "deletion_logs",
	}

	// Drop each table
	for _, table := range tables {
		query := "DROP TABLE IF EXISTS " + table + " CASCADE"
		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("❌ Failed to drop table %s: %v", table, err)
		}
	}

	log.Println("✅ All tables dropped successfully!")
}
