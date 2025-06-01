package repository

import (
	"database/sql"
	"fmt"
	"log"
	"startfront-backend/internal/domain"
	"strings"
)

// InsertScreen creates a new screen in the database
func InsertScreen(screen domain.Screen) error {
	_, err := db.DB.Exec("INSERT INTO screens (application_id, name, code, route, description, params, validate, auth_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		screen.AppID, screen.Name, screen.Code, screen.Route, screen.Description, screen.Params, screen.Validate, screen.AuthBy)

	if err != nil {
		// Check if the error is a duplicate key error
		if isDuplicateScreenError(err) {
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
func isDuplicateScreenError(err error) bool {
	if err != nil && (strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "screens_route_key") || strings.Contains(err.Error(), "screens_code_key")) {
		return true
	}
	return false
}

// GetScreensById fetches a single screen by its ID
func GetScreensById(id string) (domain.Screen, error) {
	var screen domain.Screen
	var authByName string

	err := db.DB.QueryRow("SELECT id, application_id, name, code, route, description, params, validate, auth_by FROM screens WHERE id = $1",
		id).Scan(&screen.ID, &screen.AppID, &screen.Name, &screen.Code, &screen.Route, &screen.Description, &screen.Params, &screen.Validate, &screen.AuthBy)

	// Return the error if the screen doesn't exist
	if err != nil {
		if err == sql.ErrNoRows {
			// Return custom error if the screen doesn't exist
			return screen, fmt.Errorf("screen not found")
		}
		return screen, err
	}

	// Get the usernames using the reusable function
	authByName, _ = GetUserNameByID(screen.AuthBy)

	// Assign the fetched usernames to the struct
	screen.AuthByName = authByName

	return screen, nil
}

// GetAllScreens fetches all screens from the database
func GetAllScreens() ([]domain.Screen, error) {
	rows, err := db.DB.Query("SELECT id, application_id, name, code, route, description, params, validate, auth_by FROM screens")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screens []domain.Screen
	for rows.Next() {
		var screen domain.Screen
		var authByName string

		if err := rows.Scan(&screen.ID, &screen.AppID, &screen.Name, &screen.Code, &screen.Route, &screen.Description, &screen.Params, &screen.Validate, &screen.AuthBy); err != nil {
			return nil, err
		}
		// Get the usernames using the reusable function
		authByName, _ = GetUserNameByID(screen.AuthBy)

		// Assign the fetched usernames to the struct
		screen.AuthByName = authByName
		screens = append(screens, screen)
	}
	return screens, nil
}

// UpdateScreen updates a screen in the database
func UpdateScreen(id string, screen domain.Screen) error {
	_, err := db.DB.Exec("UPDATE screens SET name = $1, code = $2, route = $3, description = $5, params = $6, validate = $7, auth_by = $8 WHERE id = $9",
		screen.Name, screen.Code, screen.Route, screen.Description, screen.Params, screen.Validate, screen.AuthBy, id)
	return err
}

// DeleteScreen deletes a screen by its code
func DeleteScreen(id string) error {
	_, err := db.DB.Exec("DELETE FROM screens WHERE id = $1", id)
	return err
}
