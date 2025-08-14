package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"startfront/models"
	"strings"
)

func HandleUsers_rolesGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		UserID int64  `json:"userId"`
		RoleID int64    `json:"roleId"`
		Status string `json:"status"`
	}

	// If there's data, unmarshal it into req
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Ensure ID is provided
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Build the query dynamically
	var m models.Users_roles
	query := db.Model(&models.Users_roles{}).Where("id = ?", req.ID)

	// Apply dynamic filters if present
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.RoleID != 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}

	if req.Status == "" {
		req.Status = "active"
	}

	// Fetch the record
	if err := query.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return the record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (get by ID or other params)",
		"data":    m,
	})
}

func HandleUsers_rolesList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `json:"page"`
		Limit  *int    `json:"limit"`
		Search *string `json:"search"`
	}

	// If there's data, unmarshal it into req
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Set default pagination
	page := 1
	limit := 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	// Build query with possible search filter
	q := db.Model(&models.Users_roles{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	// Get the total count
	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Fetch the records with pagination
	var items []models.Users_roles
	if err := q.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Prepare and send the response
	response := gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (list)",
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	if len(items) > 0 {
		response["data"] = items
	}

	c.JSON(http.StatusOK, response)
}

func HandleUsers_rolesInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Users_roles
	// Parse request data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Check if the user_id exists in the users table
	var userCount int64
	if err := db.Model(&models.Users{}).Where("id = ?", req.UserID).Count(&userCount).Error; err != nil || userCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "user_id does not exist"})
		return
	}

	// Check if the role_id exists in the roles table
	var roleCount int64
	if err := db.Model(&models.Roles{}).Where("id = ?", req.RoleID).Count(&roleCount).Error; err != nil || roleCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "role_id does not exist"})
		return
	}

	// Default status to "active" if not provided
	if req.Status == "" {
		req.Status = "active"
	}

	// Create new entry
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (inserted)",
		"data":    req,
	})
}

func HandleUsers_rolesUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		UserID int64  `json:"userId"`
		RoleID int64    `json:"roleId"`
		Status string `json:"status"`
	}

	// Parse request data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Ensure ID is provided
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Fetch record by ID
	var m models.Users_roles
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Update fields if provided
	if req.UserID != 0 {
		m.UserID = req.UserID
	}
	if req.RoleID != 0 {
		m.RoleID = req.RoleID
	}
	if req.Status != "" {
		m.Status = req.Status
	}

	// Save the updated record
	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return updated record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (updated)",
		"data":    m,
	})
}

func HandleUsers_rolesDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}

	// Parse request data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Ensure ID is provided
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Delete the record
	if err := db.Delete(&models.Users_roles{}, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (deleted)",
	})
}
