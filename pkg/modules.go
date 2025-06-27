package pkg

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Module interface {
    AutoMigrate(db *gorm.DB)
    RegisterRoutes(r *gin.Engine, db *gorm.DB)
}

var Modules []Module

func RegisterModules(r *gin.Engine, db *gorm.DB) {
    for _, m := range Modules {
        m.AutoMigrate(db)
        m.RegisterRoutes(r, db)
    }
}
