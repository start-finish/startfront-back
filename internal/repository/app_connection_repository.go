package repository

import (
	"startfront-backend/internal/domain"
)

// Insert a new App Connection
func InsertAppConnection(connection domain.AppConnection) error {
	query := `
		INSERT INTO app_connections (application_id, name, type, config)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(query, connection.ApplicationID, connection.Name, connection.Type, connection.Config)
	return err
}

// Get App Connections by Application ID
func GetAppConnectionsByApplicationID(appID int) ([]domain.AppConnection, error) {
	var connections []domain.AppConnection
	err := db.Select(&connections, `SELECT * FROM app_connections WHERE application_id = $1`, appID)
	return connections, err
}
