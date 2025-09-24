package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"startfront/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleNavigation_itemsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID    uint   `json:"id"`
		Label string `json:"label"` // was Name
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

	var m models.Navigation_items
	query := db.Model(&models.Navigation_items{}).Where("id = ?", req.ID)

	if strings.TrimSpace(req.Label) != "" {
		query = query.Where("label ILIKE ?", "%"+strings.TrimSpace(req.Label)+"%")
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
		"message": "startfront API (get by ID or other params)",
		"data":    m,
	})
}

func HandleNavigation_itemsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	q := db.Model(&models.Navigation_items{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		s := strings.TrimSpace(*req.Search)
		q = q.Where("(label ILIKE ? OR path ILIKE ?)", "%"+s+"%", "%"+s+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Navigation_items
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

func HandleNavigation_itemsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	// Unmarshal directly into your GORM model (assumes proper json tags on the model)
	var req models.Navigation_items
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// ---- Normalize & validate input ----
	// label required
	req.Label = strings.TrimSpace(req.Label)
	if req.Label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'label' is required"})
		return
	}

	// trim path; treat empty as NULL
	if req.Path != nil {
		p := strings.TrimSpace(*req.Path)
		if p == "" {
			req.Path = nil
		} else {
			req.Path = &p
		}
	}

	// menu_id must exist
	var menuCount int64
	if err := db.Model(&models.Navigation_menus{}).Where("id = ?", req.MenuID).Count(&menuCount).Error; err != nil || menuCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "menu_id does not exist"})
		return
	}

	// parent_id (if provided) must exist and belong to the same menu
	if req.ParentID != nil {
		var parent models.Navigation_items
		if err := db.First(&parent, *req.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "parent_id does not exist"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
			return
		}
		if parent.MenuID != req.MenuID {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "parent_id belongs to a different menu"})
			return
		}
	}

	// ---- Duplicate checks (context-aware) ----
	// 1) Duplicate label within (menu_id, parent_id)
	{
		var dup models.Navigation_items
		err := db.Where(
			"menu_id = ? AND COALESCE(parent_id,0) = COALESCE(?,0) AND label = ?",
			req.MenuID, req.ParentID, req.Label,
		).First(&dup).Error

		switch {
		case err == nil:
			// found duplicate
			c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate label in this menu/parent"})
			return
		case errors.Is(err, gorm.ErrRecordNotFound):
			// ok, no duplicate -> continue
		default:
			// actual DB error
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
			return
		}
	}

	// 2) Duplicate path within (menu_id) if path is set
	if req.Path != nil {
		var dupPath models.Navigation_items
		err := db.Where("menu_id = ? AND path = ?", req.MenuID, req.Path).First(&dupPath).Error

		switch {
		case err == nil:
			c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "duplicate path in this menu"})
			return
		case errors.Is(err, gorm.ErrRecordNotFound):
			// ok
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
			return
		}
	}

	// ---- Create inside a transaction ----
	var createErr error
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&req).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		createErr = err
	}
	if createErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": createErr.Error()})
		return
	}

	// ---- Reload with relations and return ----
	if err := db.Preload("Menu").Preload("Parent").First(&req, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "error loading data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (inserted)",
		"data":    req,
	})
}

func HandleNavigation_itemsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID         uint           `json:"id"`
		MenuID     *uint          `json:"menu_id"`
		ParentID   *uint          `json:"parent_id"`
		Label      *string        `json:"label"`
		Path       *string        `json:"path"`
		ScreenID   *uint          `json:"screen_id"`
		Icon       *string        `json:"icon"`
		SortOrder  *int           `json:"sort_order"`
		Properties map[string]any `json:"properties"` // if using datatypes.JSONMap switch type accordingly
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	var m models.Navigation_items
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Apply changes
	if req.MenuID != nil {
		// validate menu
		var cnt int64
		if err := db.Model(&models.Navigation_menus{}).Where("id = ?", *req.MenuID).Count(&cnt).Error; err != nil || cnt == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "menu_id does not exist"})
			return
		}
		m.MenuID = *req.MenuID
	}
	if req.ParentID != nil {
		m.ParentID = req.ParentID
	}
	if req.Label != nil {
		l := strings.TrimSpace(*req.Label)
		if l == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'label' cannot be empty"})
			return
		}
		m.Label = l
	}
	if req.Path != nil {
		p := strings.TrimSpace(*req.Path)
		if p == "" {
			m.Path = nil
		} else {
			m.Path = &p
		}
	}
	if req.ScreenID != nil {
		m.ScreenID = req.ScreenID
	}
	if req.Icon != nil {
		m.Icon = req.Icon
	}
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}
	if req.Properties != nil {
		// if your model uses datatypes.JSON (bytes), marshal the map:
		b, _ := json.Marshal(req.Properties)
		m.Properties = b
	}

	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (updated)",
		"data":    m,
	})
}

func HandleNavigation_itemsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Navigation_items{}, req.ID)
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
