package user

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/start-finish/startfront-app/internal/models"
)

type Handler struct {
    Service *Service
}

func NewHandler(s *Service) *Handler {
    return &Handler{Service: s}
}

func (h *Handler) GetUsers(c *gin.Context) {
    users, err := h.Service.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.Service.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}
