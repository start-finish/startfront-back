package auth

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "github.com/start-finish/startfront-app/internal/models"
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
    if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

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
