package auth

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    s := NewService(db)
    h := NewHandler(s)

    r.POST("/login", h.Login)
}
