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

func HandleWidget_instancesGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
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

	var m models.Widget_instances
	query := db.Model(&models.Widget_instances{}).Where("id = ?", req.ID)

	if req.Title != "" {
		query = query.Where("title ILIKE ?", "%"+req.Title+"%")
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

func HandleWidget_instancesList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	page := 1
	limit := 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	q := db.Model(&models.Widget_instances{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("title ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}
	if req.Scope != nil && *req.Scope != "" {
		q = q.Where("scope = ?", *req.Scope)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Widget_instances
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

func HandleWidget_instancesInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Widget_instances
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid request data: " + err.Error(),
		})
		return
	}

	// Validate scope rules (to match DB CHECK constraint)
	switch req.Scope {
	case "global":
		if req.ClientID != nil || req.ScreenID != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "global scope cannot have client_id or screen_id",
			})
			return
		}
	case "client":
		if req.ClientID == nil || req.ScreenID != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "client scope requires client_id and must not have screen_id",
			})
			return
		}
	case "screen":
		if req.ClientID == nil || req.ScreenID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "screen scope requires both client_id and screen_id",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid scope value (must be global, client, or screen)",
		})
		return
	}

	// Check Widget exists
	var widget models.Widgets
	if err := db.First(&widget, req.WidgetID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid widget_id",
		})
		return
	}

	// Optional: Check Client exists
	if req.ClientID != nil {
		var client models.Clients
		if err := db.First(&client, *req.ClientID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "invalid client_id",
			})
			return
		}
	}

	// Optional: Check Screen exists
	if req.ScreenID != nil {
		var screen models.Screens
		if err := db.First(&screen, *req.ScreenID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "invalid screen_id",
			})
			return
		}
	}

	// Duplicate check (title+scope+client+screen combo)
	if req.Title != nil && *req.Title != "" {
		var existing models.Widget_instances
		if err := db.Where(
			"title = ? AND scope = ? AND COALESCE(client_id,0) = COALESCE(?,0) AND COALESCE(screen_id,0) = COALESCE(?,0)",
			*req.Title, req.Scope, req.ClientID, req.ScreenID,
		).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "duplicate widget instance",
			})
			return
		}
	}

	// Create new record
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Reload with relationships for response
	if err := db.
		Preload("Widget").
		Preload("Client").
		Preload("Screen").
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

func HandleWidget_instancesUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID        uint            `json:"id"`
		Title     *string         `json:"title"`
		Scope     *string         `json:"scope"`
		ClientID  *int64          `json:"client_id"`
		ScreenID  *int64          `json:"screen_id"`
		Config    json.RawMessage `json:"config"`
		SortOrder *int            `json:"sort_order"`
		IsActive  *bool           `json:"is_active"`
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

	var m models.Widget_instances
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

	// Apply updates
	if req.Title != nil {
		m.Title = req.Title
	}
	if req.Scope != nil {
		m.Scope = *req.Scope
	}
	if req.ClientID != nil {
		m.ClientID = req.ClientID
	}
	if req.ScreenID != nil {
		m.ScreenID = req.ScreenID
	}
	if len(req.Config) > 0 {
		m.Config = datatypes.JSON(req.Config)
	}
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	}

	// Validate scope rules (to match DB CHECK constraint)
	switch m.Scope {
	case "global":
		if m.ClientID != nil || m.ScreenID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "global scope cannot have client_id or screen_id"})
			return
		}
	case "client":
		if m.ClientID == nil || m.ScreenID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "client scope requires client_id and must not have screen_id"})
			return
		}
	case "screen":
		if m.ClientID == nil || m.ScreenID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "screen scope requires both client_id and screen_id"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid scope value (must be global, client, or screen)"})
		return
	}

	// Duplicate check (exclude self)
	if m.Title != nil && *m.Title != "" {
		var existing models.Widget_instances
		if err := db.Where(
			"title = ? AND scope = ? AND COALESCE(client_id,0) = COALESCE(?,0) AND COALESCE(screen_id,0) = COALESCE(?,0) AND id <> ?",
			*m.Title, m.Scope, m.ClientID, m.ScreenID, m.ID,
		).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "duplicate widget instance",
			})
			return
		}
	}

	// Save the updated record
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

func HandleWidget_instancesDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Widget_instances{}, req.ID)
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
		"message": "widget instance deleted",
	})
}
