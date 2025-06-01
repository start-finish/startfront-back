package repository

import (
	"errors"
	"startfront-backend/internal/domain"
)

// Insert a new Application Collaborator
func InsertApplicationCollaborator(collaborator domain.ApplicationCollaborator) error {
	query := `
		INSERT INTO application_collaborators (application_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(query, collaborator.ApplicationID, collaborator.UserID, collaborator.Role)
	return err
}

// GetApplicationCollaboratorsByApplicationID fetches collaborators based on application ID
func GetApplicationCollaboratorsByApplicationID(applicationID int) ([]domain.ApplicationCollaborator, error) {
	var collaborators []domain.ApplicationCollaborator
	err := db.Select(&collaborators, `SELECT * FROM application_collaborators WHERE application_id = $1`, applicationID)
	if err != nil {
		return nil, errors.New("could not fetch collaborators")
	}
	return collaborators, nil
}

// ListApplicationCollaborators retrieves all collaborators for a specific application
func ListApplicationCollaborators(applicationID int) ([]domain.ApplicationCollaborator, error) {
	// Call repository to get collaborators based on applicationID
	collaborators, err := GetApplicationCollaboratorsByApplicationID(applicationID)
	if err != nil {
		return nil, errors.New("unable to fetch collaborators")
	}
	return collaborators, nil
}
