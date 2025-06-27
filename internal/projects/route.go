package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/start-finish/startfront-app/pkg"
	"gorm.io/gorm"
)

type ProjectModule struct{}

func (m *ProjectModule) AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Project{})
}

func (m *ProjectModule) RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	svc := &pkg.BaseService[Project]{DB: db}
	group := r.Group("/api/projects")

	pkg.RegisterCRUDRoutes(group, svc, pkg.RouteOptions{
		EnableList:   true,
		EnableGet:    true,
		EnableCreate: true,
		EnableUpdate: true,
		EnableDelete: true,
	}, []string{"name", "description"})
}
