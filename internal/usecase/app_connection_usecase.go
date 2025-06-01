package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// CreateAppConnection creates an app connection
func CreateAppConnection(connection domain.AppConnection) error {
	return repository.InsertAppConnection(connection)
}

// ListAppConnectionsByAppID gets app connections by application ID
func ListAppConnectionsByAppID(appID int) ([]domain.AppConnection, error) {
	return repository.GetAppConnectionsByApplicationID(appID)
}
