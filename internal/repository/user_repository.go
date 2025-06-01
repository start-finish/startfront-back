package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"startfront-backend/internal/domain"
	"strings"
)

// InsertUser inserts a new user into the database
func InsertUser(user domain.User) error {
	_, err := db.DB.Exec("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		// Check if the error is a duplicate key error
		if isDuplicateUserError(err) {
			log.Println("Duplicate email/name error:", err)
			return fmt.Errorf("duplicate_found")
		}
		// Log and return the error for any other case
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

// isDuplicateUserError checks if the error is due to a unique constraint violation
func isDuplicateUserError(err error) bool {
	if err != nil && (strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "users_email_key") || strings.Contains(err.Error(), "users_name_key")) {
		return true
	}
	return false
}

// GetUser fetches a user by ID
func GetUser(id string) (domain.User, error) {
	var user domain.User
	err := db.DB.QueryRow("SELECT id, name, email, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	return user, err
}

// UpdateUser updates a user's information
func UpdateUser(id string, user domain.User) error {
	_, err := db.DB.Exec("UPDATE users SET name = $1, email = $2, role = $3 WHERE id = $4", user.Name, user.Email, user.Role, id)
	return err
}

// DeleteUser deletes a user by ID
func DeleteUser(id string) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// GetUserNameByID fetches the username for a given user ID
func GetUserNameByID(userID *int) (string, error) {
	if userID == nil {
		return "", nil
	}

	var username string
	query := "SELECT name FROM users WHERE id = $1"
	err := db.DB.QueryRow(query, *userID).Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error fetching username:", err)
		return "", fmt.Errorf("failed to fetch username")
	}
	return username, nil
}

func FindUserByEmail(email string) (domain.User, error) {
	var user domain.User

	// Example with sqlx, adapt to your DB access method:
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	err := db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
