package repository

import (
	"startfront-backend/internal/domain"
)

func InsertAuthBinding(a domain.AuthBinding) error {
	query := `
		INSERT INTO auth_bindings (
			application_id, method, login_endpoint, credentials_schema, token_path, response_user_path
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query,
		a.ApplicationID, a.Method, a.LoginEndpoint, a.CredentialsSchema, a.TokenPath, a.ResponseUserPath,
	)
	return err
}

func GetAuthBindingsByAppID(appID int) ([]domain.AuthBinding, error) {
	var bindings []domain.AuthBinding
	query := `SELECT * FROM auth_bindings WHERE application_id = $1`
	err := db.Select(&bindings, query, appID)
	return bindings, err
}
