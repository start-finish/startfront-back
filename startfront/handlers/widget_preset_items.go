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

func HandleWidget_preset_itemsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}

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

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "'id' is required",
		})
		return
	}

	var m models.Widget_preset_items
	if err := db.
		Preload("Preset").
		Preload("Widget").
		First(&m, req.ID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "record not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"data":   m,
	})
}

func HandleWidget_preset_itemsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page     *int    `json:"page"`
		Limit    *int    `json:"limit"`
		Search   *string `json:"search"`
		PresetID *int64  `json:"preset_id"`
	}

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

	page, limit := 1, 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	q := db.Model(&models.Widget_preset_items{})
	if req.PresetID != nil {
		q = q.Where("preset_id = ?", *req.PresetID)
	}
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("CAST(widget_id AS TEXT) ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Widget_preset_items
	if err := q.
		Order("sort_order ASC, id ASC").
		Limit(limit).
		Offset(offset).
		Preload("Preset").
		Preload("Widget").
		Find(&items).Error; err != nil {

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

func HandleWidget_preset_itemsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Widget_preset_items
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid request data: " + err.Error(),
		})
		return
	}

	// Validate preset exists
	var preset models.Widget_presets
	if err := db.First(&preset, req.PresetID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid preset_id",
		})
		return
	}

	// Validate widget exists
	var widget models.Widgets
	if err := db.First(&widget, req.WidgetID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid widget_id",
		})
		return
	}

	// Check duplicate: preset_id + widget_id combo
	var existing models.Widget_preset_items
	if err := db.Where(
		"preset_id = ? AND widget_id = ?",
		req.PresetID, req.WidgetID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "duplicate preset/widget combination",
		})
		return
	}

	// Default config if nil
	if len(req.DefaultConfig) == 0 {
		req.DefaultConfig = datatypes.JSON([]byte(`{}`))
	}

	// Insert
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Load relations
	if err := db.
		Preload("Preset").
		Preload("Widget").
		First(&req, req.ID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "failed to load relations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"data":   req,
	})
}

func HandleWidget_preset_itemsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID            uint            `json:"id"`
		PresetID      *int64          `json:"preset_id"`
		WidgetID      *int64          `json:"widget_id"`
		DefaultConfig json.RawMessage `json:"default_config"`
		SortOrder     *int            `json:"sort_order"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid request data: " + err.Error(),
		})
		return
	}
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "'id' is required",
		})
		return
	}

	var m models.Widget_preset_items
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "record not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Update allowed fields
	if req.PresetID != nil {
		m.PresetID = *req.PresetID
	}
	if req.WidgetID != nil {
		m.WidgetID = *req.WidgetID
	}
	if len(req.DefaultConfig) > 0 {
		m.DefaultConfig = datatypes.JSON(req.DefaultConfig)
	}
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}

	// Validate preset/widget exist if changed
	if req.PresetID != nil {
		var preset models.Widget_presets
		if err := db.First(&preset, *req.PresetID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid preset_id"})
			return
		}
	}
	if req.WidgetID != nil {
		var widget models.Widgets
		if err := db.First(&widget, *req.WidgetID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid widget_id"})
			return
		}
	}

	// Check duplicate (exclude self)
	var existing models.Widget_preset_items
	if err := db.Where(
		"preset_id = ? AND widget_id = ? AND id <> ?",
		m.PresetID, m.WidgetID, m.ID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "duplicate preset/widget combination",
		})
		return
	}

	// Save
	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"data":   m,
	})
}

func HandleWidget_preset_itemsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Widget_preset_items{}, req.ID)
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
		"message": "widget preset item deleted",
	})
}
