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

func HandleRole_permissionsGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID           uint   `json:"id"`
		RoleID       int64  `json:"role_id"`
		PermissionID int64  `json:"permission_id"`
		Status       string `json:"status"`
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
	var m models.Role_permissions
	query := db.Model(&models.Role_permissions{}).Where("id = ?", req.ID)

	// Apply dynamic filters (Name, Status, etc.)
	if req.RoleID != 0 {
		query = query.Where("role_id ILIKE ?", req.RoleID)
	}
	if req.PermissionID != 0 {
		query = query.Where("permission_id ILIKE ?", req.PermissionID)
	}

	query = query.Preload("Role").Preload("Permission")

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

func HandleRole_permissionsList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	q := db.Model(&models.Role_permissions{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.Role_permissions
	if err := q.Preload("Role").Preload("Permission").Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
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

func HandleRole_permissionsInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.Role_permissions
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	var roleCount int64
	if err := db.Model(&models.Roles{}).Where("id = ?", req.RoleID).Count(&roleCount).Error; err != nil || roleCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "role_id does not exist"})
		return
	}

	var permisisonCount int64
	if err := db.Model(&models.Permissions{}).Where("id = ?", req.PermissionID).Count(&permisisonCount).Error; err != nil || permisisonCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "permission_id does not exist"})
		return
	}

	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var result models.Role_permissions
	if err := db.Preload("Role").Preload("Permission").First(&result, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": "error preloading related data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (inserted)",
		"data":    result,
	})
}

func HandleRole_permissionsUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID           uint   `json:"id"`
		RoleID       int64  `json:"role_id"`
		PermissionID int64  `json:"permission_id"`
		Status       string `json:"status"`
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
	var m models.Role_permissions
	if err := db.First(&m, "id = ?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Validate FKs only if changing them
	if req.RoleID != 0 {
		var cnt int64
		if err := db.Model(&models.Roles{}).Where("id = ?", req.RoleID).Count(&cnt).Error; err != nil || cnt == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "role_id does not exist"})
			return
		}
	}
	if req.PermissionID != 0 {
		var cnt int64
		if err := db.Model(&models.Permissions{}).Where("id = ?", req.PermissionID).Count(&cnt).Error; err != nil || cnt == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "permission_id does not exist"})
			return
		}
	}

	// Uniqueness check for (user_id, role_id) pair excluding this row
	newRoleID := m.RoleID
	if req.RoleID != 0 {
		newRoleID = req.RoleID
	}
	newPermisisonID := m.PermissionID
	if req.PermissionID != 0 {
		newPermisisonID = req.PermissionID
	}
	var dup int64
	if err := db.Model(&models.Role_permissions{}).
		Where("role_id = ? AND permission_id = ? AND id <> ?", newRoleID, newPermisisonID, req.ID).
		Count(&dup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}
	if dup > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": "1", "status": "error", "error": "user_id and role_id pair already exists"})
		return
	}

	// Build update map
	updates := map[string]interface{}{}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}
	if req.PermissionID != 0 {
		updates["permission_id"] = req.PermissionID
	}

	// Nothing to change
	if len(updates) == 0 {
		if err := db.Preload("Role").Preload("Permission").First(&m, "id = ?", req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "0", "status": "Success", "message": "startfront API (updated - no changes)", "data": m})
		return
	}

	// Plain UPDATE; check RowsAffected
	res := db.Model(&models.Role_permissions{}).Where("id = ?", req.ID).Updates(updates)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": res.Error.Error()})
		return
	}
	if res.RowsAffected == 0 {
		// Someone deleted it between read and update, or wrong id
		c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found during update"})
		return
	}

	// Reload with relations
	var out models.Role_permissions
	if err := db.Preload("Role").Preload("Permission").First(&out, "id = ?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found after update"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return the updated record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (updated)",
		"data":    out,
	})
}

func HandleRole_permissionsDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.Role_permissions{}, req.ID)
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
