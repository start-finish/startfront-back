package handler

import (
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(au usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: au}
}

type loginRequest struct {
	Username string `json:"username"` // can be username or email
	Password   string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			response.Error(c, "User not found")
		case usecase.ErrInvalidCredential:
			response.Error(c, "Invalid password")
		default:
			response.Error(c, "Login failed")
		}
		return
	}
	response.Success(c, "Login successfully", nil)
}
