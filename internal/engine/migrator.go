package engine

import (
	"gorm.io/gorm"
)

func AutoMigrateSchema(db *gorm.DB, schemaDef Schema) error {
	model, err := BuildGormModel(schemaDef)
	if err != nil {
		return err
	}

	if err := db.Table(schemaDef.Model).AutoMigrate(model); err != nil {
		return err
	}

	return nil
}
