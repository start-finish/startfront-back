package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// CreateApplicationCollaborator creates a new application collaborator
func CreateApplicationCollaborator(collaborator domain.ApplicationCollaborator) error {
	return repository.InsertApplicationCollaborator(collaborator)
}

// ListApplicationCollaborators retrieves all collaborators for a given application
func ListApplicationCollaborators(applicationID int) ([]domain.ApplicationCollaborator, error) {
	return repository.GetApplicationCollaboratorsByApplicationID(applicationID)
}
