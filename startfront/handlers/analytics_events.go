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

func HandleAnalytics_eventsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID        uint   `json:"id"`
		EventName string `json:"event_name"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "1", "status": "error", "error": "invalid request data: " + err.Error(),
			})
			return
		}
	}

	if req.ID == 0 && req.EventName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' or 'event_name' is required"})
		return
	}

	var m models.Analytics_events
	query := db.Model(&models.Analytics_events{}).
		Preload("User").Preload("Client").Preload("Screen").Preload("WidgetInstance")

	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.EventName != "" {
		query = query.Where("event_name = ?", req.EventName)
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

func HandleAnalytics_eventsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page      *int    `json:"page"`
		Limit     *int    `json:"limit"`
		EventName *string `json:"event_name"`
		ClientID  *int64  `json:"client_id"`
		UserID    *int64  `json:"user_id"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "1", "status": "error", "error": "invalid request data: " + err.Error(),
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

	q := db.Model(&models.Analytics_events{}).
		Preload("User").Preload("Client").Preload("Screen").Preload("WidgetInstance")

	if req.EventName != nil && strings.TrimSpace(*req.EventName) != "" {
		q = q.Where("event_name ILIKE ?", "%"+strings.TrimSpace(*req.EventName)+"%")
	}
	if req.ClientID != nil {
		q = q.Where("client_id = ?", *req.ClientID)
	}
	if req.UserID != nil {
		q = q.Where("user_id = ?", *req.UserID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Analytics_events
	if err := q.Order("occurred_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
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

func HandleAnalytics_eventsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Analytics_events
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "invalid request data: " + err.Error(),
		})
		return
	}

	// Required field check
	if req.EventName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "'event_name' is required",
		})
		return
	}

	// Foreign key validation
	if req.ClientID != nil {
		var client models.Clients
		if err := db.First(&client, *req.ClientID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid client_id"})
			return
		}
	}
	if req.UserID != nil {
		var user models.Users
		if err := db.First(&user, *req.UserID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid user_id"})
			return
		}
	}
	if req.ScreenID != nil {
		var screen models.Screens
		if err := db.First(&screen, *req.ScreenID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid screen_id"})
			return
		}
	}
	if req.WidgetInstanceID != nil {
		var widget models.Widget_instances
		if err := db.First(&widget, *req.WidgetInstanceID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid widget_instance_id"})
			return
		}
	}

	// Ensure meta JSON is initialized
	if len(req.Meta) == 0 {
		req.Meta = datatypes.JSON([]byte(`{}`))
	}

	// Insert event
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "1", "status": "error", "error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0", "status": "Success", "data": req,
	})
}

func HandleAnalytics_eventsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID        uint            `json:"id"`
		EventName *string         `json:"event_name"`
		Meta      json.RawMessage `json:"meta"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "invalid request data: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "'id' is required",
		})
		return
	}

	var m models.Analytics_events
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	if req.EventName != nil {
		m.EventName = *req.EventName
	}
	if len(req.Meta) > 0 {
		m.Meta = datatypes.JSON(req.Meta)
	}

	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": m})
}

func HandleAnalytics_eventsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "invalid request data: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "1", "status": "error", "error": "'id' is required",
		})
		return
	}

	res := db.Delete(&models.Analytics_events{}, req.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": res.Error.Error()})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0", "status": "Success", "message": "event deleted",
	})
}
