package handlers

import (
	"encoding/json"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func HandleSettingsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	if req.ID == 0 && req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' or 'key' is required"})
		return
	}

	var m models.Settings
	query := db.Model(&models.Settings{}).Preload("Client")

	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.Key != "" {
		query = query.Where("key = ?", req.Key)
	}

	if err := query.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": m})
}

func HandleSettingsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	q := db.Model(&models.Settings{}).Preload("Client")
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("key ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}
	if req.Scope != nil && *req.Scope != "" {
		q = q.Where("scope = ?", *req.Scope)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Settings
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

func HandleSettingsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Settings
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Validate scope rules
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

	// Duplicate check
	var existing models.Settings
	if err := db.Where(
		"key = ? AND scope = ? AND COALESCE(client_id, 0) = COALESCE(?, 0)",
		req.Key, req.Scope, req.ClientID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate setting"})
		return
	}

	if len(req.Value) == 0 {
		req.Value = datatypes.JSON([]byte(`{}`))
	}

	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": req})
}

func HandleSettingsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID       uint            `json:"id"`
		Key      *string         `json:"key"`
		Scope    *string         `json:"scope"`
		ClientID *int64          `json:"client_id"`
		Value    json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	var m models.Settings
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	if req.Key != nil {
		m.Key = *req.Key
	}
	if req.Scope != nil {
		m.Scope = *req.Scope
	}
	if req.ClientID != nil {
		m.ClientID = req.ClientID
	}
	if len(req.Value) > 0 {
		m.Value = datatypes.JSON(req.Value)
	}

	// Validate scope
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
	var existing models.Settings
	if err := db.Where(
		"key = ? AND scope = ? AND COALESCE(client_id, 0) = COALESCE(?, 0) AND id <> ?",
		m.Key, m.Scope, m.ClientID, m.ID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate setting"})
		return
	}

	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": m})
}

func HandleSettingsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Settings{}, req.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": res.Error.Error()})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "message": "setting deleted"})
}
