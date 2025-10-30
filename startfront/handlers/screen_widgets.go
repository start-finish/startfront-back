package handlers

import (
	"encoding/json"
	"net/http"
	"startfront/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Get single screen_widget by composite key
func HandleScreen_widgetsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint64 `json:"id"` // Use 'id' as the primary key
	}

	// Parse the JSON data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "id is required"})
		return
	}

	var item models.Screen_widgets
	if err := db.Preload("Screen").
		Preload("WidgetInstance").
		Where("id = ?", req.ID).
		First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": item})
}

// List screen_widgets with pagination
func HandleScreen_widgetsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page  *int `json:"page"`
		Limit *int `json:"limit"`
	}

	// Parse the JSON data
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": err.Error()})
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

	var total int64
	if err := db.Model(&models.Screen_widgets{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Screen_widgets
	if err := db.Preload("Screen").
		Preload("WidgetInstance").
		Order("screen_id ASC, widget_instance_id ASC").
		Limit(limit).Offset(offset).
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

// Insert new screen_widget
func HandleScreen_widgetsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Screen_widgets
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "invalid request data: " + err.Error(),
		})
		return
	}

	// Validate relationships
	var screen models.Screens
	if err := db.First(&screen, req.ScreenID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid screen_id"})
		return
	}
	var widgetInstance models.Widget_instances
	if err := db.First(&widgetInstance, req.WidgetInstanceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid widget_instance_id"})
		return
	}

	// Create record
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return record with preloaded relations
	if err := db.Preload("Screen").
		Preload("WidgetInstance").
		First(&req, "screen_id = ? AND widget_instance_id = ?", req.ScreenID, req.WidgetInstanceID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to load created record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": req})
}

// Update screen_widget
func HandleScreen_widgetsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID                 uint64                 `json:"id"`                         // Use 'id' for update
		NewScreenID         *uint64                `json:"new_screen_id"`
		NewWidgetInstanceID *uint64                `json:"new_widget_instance_id"`
		Layout              map[string]interface{} `json:"layout"`
		SortOrder           *int                   `json:"sort_order"`
	}

	// Parse the JSON data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid JSON format"})
		return
	}

	// Validate that id is provided
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "id is required"})
		return
	}

	// Fetch the existing record based on ID
	var m models.Screen_widgets
	if err := db.Where("id = ?", req.ID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// If the screen ID was updated, ensure that the new screen exists
	if req.NewScreenID != nil && *req.NewScreenID != m.ScreenID {
		var updatedScreen models.Screens
		if err := db.First(&updatedScreen, *req.NewScreenID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "screen not found"})
			return
		}
		m.ScreenID = uint64(updatedScreen.ID) // Update the screen ID
	}

	// If the widget instance ID was updated, ensure that the new widget instance exists
	if req.NewWidgetInstanceID != nil && *req.NewWidgetInstanceID != m.WidgetInstanceID {
		var updatedWidgetInstance models.Widget_instances
		if err := db.First(&updatedWidgetInstance, *req.NewWidgetInstanceID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "widget instance not found"})
			return
		}
		m.WidgetInstanceID = uint64(updatedWidgetInstance.ID) // Update the widget instance ID
	}

	// Update layout if provided
	if req.Layout != nil {
		layoutJSON, err := json.Marshal(req.Layout)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to marshal layout"})
			return
		}
		m.Layout = datatypes.JSON(layoutJSON)
	}

	// Update sort order if provided
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}

	// Perform the update
	result := db.Model(&m).Updates(models.Screen_widgets{
		ScreenID:         m.ScreenID,
		WidgetInstanceID: m.WidgetInstanceID,
		Layout:           m.Layout,
		SortOrder:        m.SortOrder,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to update record: " + result.Error.Error()})
		return
	}

	// Check if any rows were affected (updated)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "no records were updated"})
		return
	}

	// Return the updated record
	if err := db.Preload("Screen").Preload("WidgetInstance").
		First(&m, "id = ?", m.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "failed to load updated record"})
		return
	}

	// Return the updated record with relationships
	c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "data": m})
}


// Delete screen_widget
func HandleScreen_widgetsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint64 `json:"id"` // Use 'id' for delete
	}

	// Parse the JSON data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "id is required"})
		return
	}

	// Fetch the record before deleting
	var deleted models.Screen_widgets
	if err := db.Preload("Screen").
		Preload("WidgetInstance").
		Where("id = ?", req.ID).
		First(&deleted).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	if err := db.Delete(&deleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "deleted successfully",
		"data":    deleted, // Optional: return deleted record
	})
}
