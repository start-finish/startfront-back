package repository

import (
	"startfront-backend/internal/domain"
)

// // Insert a new Widget Preset
// func InsertWidgetPreset(preset domain.WidgetPreset) error {
// 	query := `
// 		INSERT INTO widget_presets (name, type, props, auth_by)
// 		VALUES ($1, $2, $3, $4)
// 	`
// 	_, err := db.Exec(query, preset.Name, preset.Type, preset.Props, preset.AuthBy)
// 	return err
// }

// Get all Widget Presets
func GetAllWidgetPresets() ([]domain.WidgetPreset, error) {
	var presets []domain.WidgetPreset
	err := db.Select(&presets, `SELECT * FROM widget_presets`)
	return presets, err
}
