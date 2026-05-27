package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleScreensGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		RoutePath   string `json:"route_path"`
		Description string `json:"description"`
		IsActive    string `json:"is_active" gorm:"default:'0'"`
		CreatedBy   string `json:"created_by"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// ID is required
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Build query dynamically based on available parameters
	var m models.Screens
	query := db.Model(&models.Screens{}).Where("id = ?", req.ID)

	// Apply dynamic filters (Name, Status, etc.)
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.RoutePath != "" {
		query = query.Where("route_path ILIKE ?", "%"+req.RoutePath+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.IsActive != "" {
		query = query.Where("is_active ILIKE ?", "%"+req.IsActive+"%")
	}
	if req.CreatedBy != "" {
		query = query.Where("created_by ILIKE ?", "%"+req.CreatedBy+"%")
	}

	// Fetch the record based on ID and optional filters
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

func HandleScreensList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	q := db.Model(&models.Screens{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Screens
	if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
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

func HandleScreensInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Screens
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Use reflection to check for duplicates dynamically
	requiredFields := []struct {
		FieldName  string
		FieldValue interface{}
	}{
		{"name", req.Name},
		{"route_path", req.RoutePath},
		// Add more required fields here, e.g., {"FieldName", req.FieldValue}
	}

	// Check for duplicates on each required field
	for _, field := range requiredFields {
		if field.FieldValue == "" {
			// Skip empty fields
			continue
		}

		var existingRecord models.Screens
		if err := db.Where(fmt.Sprintf("%s = ?", field.FieldName), field.FieldValue).First(&existingRecord).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "1",
				"status":  "error",
				"error":   fmt.Sprintf("duplicate field: %s", field.FieldName),
			})
			return
		}
	}

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

func HandleScreensUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		RoutePath   string `json:"route_path"`
		Description string `json:"description"`
		IsActive    string `json:"is_active" gorm:"default:'0'"`
		CreatedBy   string `json:"created_by"`
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
	var m models.Screens
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.RoutePath != "" {
		m.RoutePath = req.RoutePath
	}
	if req.Description != "" {
		m.Description = req.Description
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.IsActive != "" {
		m.IsActive = req.IsActive
	}
	if req.CreatedBy != "" {
		m.CreatedBy = req.CreatedBy
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

func HandleScreensDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Screens{}, req.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": res.Error.Error()})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (deleted)",
	})
}
