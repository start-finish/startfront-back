package handlers

import (
	"encoding/json"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hashes the Users's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword compares the provided password with the stored hash
func ComparePassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

func HandleLogin(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Fetch user by username
	var m models.Users
	if err := db.Where("username = ?", req.Username).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Compare the password
	if !ComparePassword(m.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "1", "status": "error", "error": "invalid password"})
		return
	}

	// Return Users info
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "login successful",
		// "data":    m,
	})
}

func HandleSignUp(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Users
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Hash the password before saving to the database
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to hash password: " + err.Error()})
		return
	}
	req.Password = hashedPassword

	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "Users created successfully",
		"data":    req,
	})
}

func HandleUsersList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `json:"page"`
		Limit  *int    `json:"limit"`
		Search *string `json:"search"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	page := 1
	limit := 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	q := db.Model(&models.Users{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("username ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Users
	if err := q.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	response := gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (list)",
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	}
	if len(items) > 0 {
		response["data"] = items
	}

	c.JSON(http.StatusOK, response)
}

func HandleUsersUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID       uint   `json:"id"`
		Username     string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// ID is required
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Fetch the record by ID
	var m models.Users
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to hash password: " + err.Error()})
			return
		}
		m.Password = hashedPassword
	}

	// Update other fields if provided
	if req.Username != "" {
		m.Username = req.Username
	}

	// Save the updated record
	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return the updated record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (updated)",
		"data":    m,
	})
}

func HandleUsersDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	if err := db.Delete(&models.Users{}, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (deleted)",
	})
}
