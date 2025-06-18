package user

import (
    "github.com/start-finish/startfront-app/internal/models"
    "gorm.io/gorm"
)

type Service struct {
    DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{DB: db}
}

func (s *Service) GetAllUsers() ([]models.User, error) {
    var users []models.User
    result := s.DB.Find(&users)
    return users, result.Error
}

func (s *Service) CreateUser(user *models.User) error {
    return s.DB.Create(user).Error
}
