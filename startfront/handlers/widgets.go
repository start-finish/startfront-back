package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func HandleWidgetsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID  uint   `json:"id"`
		Key string `json:"key"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Either ID or Key is required
	if req.ID == 0 && strings.TrimSpace(req.Key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' or 'key' is required"})
		return
	}

	var m models.Widgets
	query := db.Model(&models.Widgets{})
	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.Key != "" {
		query = query.Where("key = ?", strings.TrimSpace(req.Key))
	}

	if err := query.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "Widget retrieved successfully",
		"data":    m,
	})
}

// ---------- LIST ----------
func HandleWidgetsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `json:"page"`
		Limit  *int    `json:"limit"`
		Search *string `json:"search"`
	}

	// Parse JSON input
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "invalid request data: " + err.Error(),
			})
			return
		}
	}

	// Default pagination
	page := 1
	limit := 10

	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}

	// Search query
	q := db.Model(&models.Widgets{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("key ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	// Count total items
	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// ---- Calculate totalPages ----
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	if totalPages == 0 {
		totalPages = 1
	}

	// ---- Validate page ----
	if page > totalPages {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  fmt.Sprintf("page %d exceeds total pages %d", page, totalPages),
		})
		return
	}

	// Fetch paginated items
	offset := (page - 1) * limit
	var items []models.Widgets
	if err := q.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Count widget categories
	type CountResult struct {
		Total   int `json:"total"`
		Control int `json:"control"`
		Input   int `json:"input"`
		Layout  int `json:"layout"`
	}

	var counts CountResult
	db.Raw(`
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE category = 'control') AS control,
			COUNT(*) FILTER (WHERE category = 'input') AS input,
			COUNT(*) FILTER (WHERE category = 'layout') AS layout
		FROM widgets
	`).Scan(&counts)

	// Build final response
	response := gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "Widgets list",
		"data": gin.H{
			"counts":  counts,
			"widgets": items,
		},
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ---------- INSERT ----------
func HandleWidgetsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Widgets

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "invalid request data: " + err.Error(),
		})
		return
	}

	// -------- VALIDATION --------

	if strings.TrimSpace(req.Key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'key' is required"})
		return
	}

	if strings.TrimSpace(req.Label) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'label' is required"})
		return
	}

	if strings.TrimSpace(req.Category) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'category' is required"})
		return
	}

	// Verify built-in vs. version rules
	if req.IsBuiltin && req.Version != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "built-in widgets cannot have a version",
		})
		return
	}

	if !req.IsBuiltin && req.Version == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "custom widgets must have a version",
		})
		return
	}

	// Check duplicate key
	var count int64
	db.Model(&models.Widgets{}).Where("key = ?", req.Key).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"code": "1", "status": "error",
			"error": "duplicate key",
		})
		return
	}

	// Insert
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "1", "status": "error", "error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0", "status": "Success",
		"message": "Widget inserted",
		"data":    req,
	})
}

// ---------- UPDATE ----------
func HandleWidgetsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID           uint            `json:"id"`
		Key          *string         `json:"key"`
		Label        *string         `json:"label"`
		Category     *string         `json:"category"`
		IconType     *string         `json:"icon_type"`
		IconValue    *string         `json:"icon_value"`
		IsBuiltin    *bool           `json:"is_builtin"`
		Version      *string         `json:"version"`
		ConfigSchema json.RawMessage `json:"config_schema"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "invalid request data: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	var m models.Widgets
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// -------- UPDATE FIELDS --------

	if req.Key != nil {
		// Check duplicate key if changed
		var count int64
		db.Model(&models.Widgets{}).Where("key = ? AND id != ?", *req.Key, req.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"code": "1", "status": "error", "error": "duplicate key",
			})
			return
		}
		m.Key = *req.Key
	}

	if req.Label != nil {
		m.Label = *req.Label
	}

	if req.Category != nil {
		m.Category = *req.Category
	}

	if req.IconType != nil {
		m.IconType = req.IconType
	}

	if req.IconValue != nil {
		m.IconValue = req.IconValue
	}

	if req.IsBuiltin != nil {
		m.IsBuiltin = *req.IsBuiltin
	}

	if req.Version != nil {
		m.Version = req.Version
	}

	// Validate built-in rules again after update
	if m.IsBuiltin && m.Version != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "built-in widgets cannot have a version",
		})
		return
	}

	if !m.IsBuiltin && m.Version == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error",
			"error": "custom widgets must have a version",
		})
		return
	}

	// Config schema update
	if len(req.ConfigSchema) > 0 {
		m.ConfigSchema = datatypes.JSON(req.ConfigSchema)
	}

	// Save updates
	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0", "status": "Success",
		"message": "Widget updated",
		"data":    m,
	})
}

// ---------- DELETE ----------
func HandleWidgetsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Widgets{}, req.ID)
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
		"message": "Widget deleted",
	})
}
