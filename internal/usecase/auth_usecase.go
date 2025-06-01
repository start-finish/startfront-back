package usecase

import (
	"errors"
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type AuthUsecase interface {
	Login(email, password string) (*domain.User, error)
}

type authUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (u *authUsecase) Login(identifier, password string) (*domain.User, error) {
	user, err := u.repo.GetUserByIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
