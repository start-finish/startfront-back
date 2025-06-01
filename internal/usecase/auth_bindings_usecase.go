package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

func CreateAuthBinding(a domain.AuthBinding) error {
	return repository.InsertAuthBinding(a)
}

func ListAuthBindingsByAppID(appID int) ([]domain.AuthBinding, error) {
	return repository.GetAuthBindingsByAppID(appID)
}
