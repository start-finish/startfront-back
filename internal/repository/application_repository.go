package repository

import (
	"fmt"
	"log"
	"startfront-backend/internal/domain"
	"strings"
)

// InsertApplication creates a new application in the database
func InsertApplication(application domain.Application) error {
	_, err := db.DB.Exec("INSERT INTO applications (name, code, route, description, auth_by) VALUES ($1, $2, $3, $4, $5)",
		application.Name, application.Code, application.Route, application.Description, application.AuthBy)

	if err != nil {
		// Check if the error is a duplicate key error
		if isDuplicateAppError(err) {
			log.Println("Duplicate code or route error:", err)
			return fmt.Errorf("duplicate_found")
		}
		// Log and return the error for any other case
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

// isDuplicateKeyError checks if the error is due to a unique constraint violation
func isDuplicateAppError(err error) bool {
	if err != nil && (strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "applications_route_key") || strings.Contains(err.Error(), "applications_code_key")) {
		return true
	}
	return false
}

// GetApplicationById fetches an application by its code
func GetApplicationById(id string) (domain.Application, error) {
	var app domain.Application
	var authByName string

	err := db.DB.QueryRow("SELECT id, name, code, route, description, auth_by FROM applications WHERE id = $1",
		id).Scan(&app.ID, &app.Name, &app.Code, &app.Route, &app.Description, &app.AuthBy)

	// If there's an error, return it
	if err != nil {
		return app, err
	}

	// Get the usernames using the reusable function
	authByName, _ = GetUserNameByID(app.AuthBy)

	// Assign the fetched usernames to the struct
	app.AuthByName = authByName

	return app, nil

}

// GetAllApplications fetches all applications
func GetAllApplications() ([]domain.Application, error) {
	rows, err := db.DB.Query("SELECT id, name, code, route, description, auth_by FROM applications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []domain.Application
	for rows.Next() {
		var app domain.Application
		var authByName string

		if err := rows.Scan(&app.ID, &app.Name, &app.Code, &app.Route, &app.Description, &app.AuthBy); err != nil {
			return nil, err
		}

		// Get the usernames using the reusable function
		authByName, _ = GetUserNameByID(app.AuthBy)

		// Assign the fetched usernames to the struct
		app.AuthByName = authByName
		apps = append(apps, app)
	}
	return apps, nil
}

// UpdateApplication updates an application's details
func UpdateApplication(id string, application domain.Application) error {
	_, err := db.DB.Exec("UPDATE applications SET name = $1, code = $2, route = $3, description = $4, auth_by = $5 WHERE id = $6",
		application.Name, application.Code, application.Route, application.Description, application.AuthBy, id)
	return err
}

// DeleteApplication deletes an application
func DeleteApplication(id string) error {
	_, err := db.DB.Exec("DELETE FROM applications WHERE id = $1", id)
	return err
}
