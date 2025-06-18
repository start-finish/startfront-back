package user

import (
	"github.com/gin-gonic/gin"
	"github.com/start-finish/startfront-app/internal/auth"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	s := NewService(db)
	h := NewHandler(s)

	userGroup := r.Group("/users")
	userGroup.Use(auth.JWTMiddleware()) // ✅ Protect all /users routes
	{
		userGroup.GET("/", auth.RequireRole("admin"), h.GetUsers) // Only admins
		userGroup.POST("/", h.CreateUser)                         // All authenticated users
	}
}
