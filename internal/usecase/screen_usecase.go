package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// CreateScreen creates a new screen
func CreateScreen(screen domain.Screen) error {
	return repository.InsertScreen(screen)
}

// GetScreensById fetches screens by application code
func GetScreensById(id string) (domain.Screen, error) {
	return repository.GetScreensById(id)
}

// ListScreens lists all Screens
func ListScreens() ([]domain.Screen, error) {
	return repository.GetAllScreens()
}

// UpdateScreen updates screen details
func UpdateScreen(id string, screen domain.Screen) error {
	return repository.UpdateScreen(id, screen)
}

// DeleteScreen deletes a screen
func DeleteScreen(id string) error {
	return repository.DeleteScreen(id)
}
