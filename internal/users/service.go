package users

import (
	"github.com/start-finish/startfront-app/models"
	"github.com/start-finish/startfront-app/pkg"
	"gorm.io/gorm"
)

type UserService struct {
    pkg.BaseService[models.Users]
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        BaseService: pkg.BaseService[models.Users]{DB: db},
    }
}
