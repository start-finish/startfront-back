package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// CreateApplication creates a new application
func CreateApplication(application domain.Application) error {
	return repository.InsertApplication(application)
}

// GetApplication fetches application details by code
func GetApplication(id string) (domain.Application, error) {
	return repository.GetApplicationById(id)
}

// ListApplications lists all applications
func ListApplications() ([]domain.Application, error) {
	return repository.GetAllApplications()
}

// UpdateApplication updates application details
func UpdateApplication(id string, application domain.Application) error {
	return repository.UpdateApplication(id, application)
}

// DeleteApplication deletes an application
func DeleteApplication(id string) error {
	return repository.DeleteApplication(id)
}
