package users

import (
	"github.com/gin-gonic/gin"
	"github.com/start-finish/startfront-app/models"
	"github.com/start-finish/startfront-app/pkg"
	"gorm.io/gorm"
)

type UserModule struct{}

func (m *UserModule) AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Users{})
}

func (m *UserModule) RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	svc := &pkg.BaseService[models.Users]{DB: db}

	pkg.RegisterUnifiedRoute(
		r,
		svc,
		pkg.RouteOptions{
			EnableList: true,
			ListMsgID:  "0",

			EnableGet: true,
			GetMsgID:  "03",

			EnableCreate: true,
			CreateMsgID:  "01",

			EnableUpdate: true,
			UpdateMsgID:  "04",

			EnableDelete: true,
			DeleteMsgID:  "05",
		},
		[]string{"email", "username"},
		nil, // ✅ using default handler, not custom one
		&pkg.RequestTypes{
			Create: &UserCreatePayload{},
			Update: &UserUpdatePayload{},
			Get:    &UserIDPayload{},
			Delete: &UserIDPayload{},
			List:   &UserListPayload{},
		},
	)
}

// --- Request Structs ---

type UserCreatePayload struct {
	Email    string                 `json:"email" binding:"required,email"`
	Username string                 `json:"username" binding:"required"`
	Password string                 `json:"password" binding:"required"`
	Meta     map[string]interface{} `json:"meta"`
}

type UserUpdatePayload struct {
	ID       uint   `json:"id" binding:"required"`
	FullName string `json:"full_name"`
}

type UserIDPayload struct {
	ID uint `json:"id" binding:"required"`
}

type UserListPayload struct {
	Page  int    `form:"page"`
	Email string `form:"email"`
	Role  string `form:"role"`
}
