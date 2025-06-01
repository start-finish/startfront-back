package repository

import (
	"database/sql"
	"errors"
	"startfront-backend/internal/domain"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByIdentifier(identifier string) (*domain.User, error)
}

type authRepository struct {
	db *sql.DB
}

// GetUserByEmail implements AuthRepository (you can implement this later or remove if unused)
func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
	panic("unimplemented")
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// This is the correct method to keep:
func (r *authRepository) GetUserByIdentifier(identifier string) (*domain.User, error) {
	user := domain.User{}
	query := `SELECT id, name, email, password, role FROM users WHERE LOWER(email) = LOWER($1) OR LOWER(name) = LOWER($1)`
	err := r.db.QueryRow(query, identifier).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
