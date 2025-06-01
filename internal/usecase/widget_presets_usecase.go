package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// // Create a Widget Preset
// func CreateWidgetPreset(preset domain.WidgetPreset) error {
// 	// Add custom logic here if needed (e.g., validations)
// 	return repository.InsertWidgetPreset(preset)
// }

// Get Widget Presets
func GetWidgetPresets() ([]domain.WidgetPreset, error) {
	return repository.GetAllWidgetPresets()
}
