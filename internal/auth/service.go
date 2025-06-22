package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/start-finish/startfront-app/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
    DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{DB: db}
}

func (s *Service) Authenticate(email, password string) (*models.User, error) {
	var user models.User

	// Find user by email (or you can use username if needed)
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare hashed password with input password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password does not match
		return nil, fmt.Errorf("invalid credentials")
	}

	// Password matches — return user
	return &user, nil
}

func (s *Service) GenerateJWT(user *models.User) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "default_secret"
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString([]byte(secret))
}
