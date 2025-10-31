package handlers

import (
	"encoding/json"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleWidget_presetsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	var m models.Widget_presets
	query := db.Model(&models.Widget_presets{}).Where("id = ?", req.ID)

	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
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
		"code":   "0",
		"status": "Success",
		"data":   m,
	})
}

func HandleWidget_presetsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `json:"page"`
		Limit  *int    `json:"limit"`
		Search *string `json:"search"`
		Scope  *string `json:"scope"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	page, limit := 1, 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	q := db.Model(&models.Widget_presets{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}
	if req.Scope != nil && *req.Scope != "" {
		q = q.Where("scope = ?", *req.Scope)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Widget_presets
	if err := q.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"meta":   gin.H{"page": page, "limit": limit, "total": total},
		"data":   items,
	})
}

func HandleWidget_presetsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Widget_presets
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Validate scope
	switch req.Scope {
	case "global":
		if req.ClientID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "global scope cannot have client_id"})
			return
		}
	case "client":
		if req.ClientID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "client scope requires client_id"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid scope value (must be global or client)"})
		return
	}

	// Validate client existence if provided
	if req.ClientID != nil {
		var client models.Clients
		if err := db.First(&client, *req.ClientID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid client_id"})
			return
		}
	}

	// Check duplicate (name + scope + client_id)
	var existing models.Widget_presets
	if err := db.Where(
		"name = ? AND scope = ? AND COALESCE(client_id,0) = COALESCE(?,0)",
		req.Name, req.Scope, req.ClientID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate widget preset"})
		return
	}

	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": req})
}

func HandleWidget_presetsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID          uint    `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Scope       *string `json:"scope"`
		ClientID    *int64  `json:"client_id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	var m models.Widget_presets
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Apply updates
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.Scope != nil {
		m.Scope = *req.Scope
	}
	if req.ClientID != nil {
		m.ClientID = req.ClientID
	}

	// Validate updated scope
	switch m.Scope {
	case "global":
		if m.ClientID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "global scope cannot have client_id"})
			return
		}
	case "client":
		if m.ClientID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "client scope requires client_id"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid scope value (must be global or client)"})
		return
	}

	// Duplicate check (exclude itself)
	if m.Name != "" {
		var existing models.Widget_presets
		if err := db.Where(
			"name = ? AND scope = ? AND COALESCE(client_id,0) = COALESCE(?,0) AND id <> ?",
			m.Name, m.Scope, m.ClientID, m.ID,
		).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate widget preset"})
			return
		}
	}

	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": m})
}

func HandleWidget_presetsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Widget_presets{}, req.ID)
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
		"message": "widget preset deleted",
	})
}
