package projects

import (
	"github.com/start-finish/startfront-app/pkg"
	"gorm.io/gorm"
)

type ProjectService struct {
    pkg.BaseService[Project]
}

func NewProjectService(db *gorm.DB) *ProjectService {
    return &ProjectService{
        BaseService: pkg.BaseService[Project]{DB: db},
    }
}
