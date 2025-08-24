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

// Handle GET by ID or filters
func HandleUser_rolesGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		UserID int64  `json:"user_id"`
		RoleID int64  `json:"role_id"`
	}

	// If there's data, unmarshal it into req
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Ensure ID is provided
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Build the query dynamically
	var m models.User_roles
	query := db.Model(&models.User_roles{}).Where("id = ?", req.ID)

	// Apply dynamic filters if present
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.RoleID != 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}

	// Preload related user and role data
	query = query.Preload("User").Preload("Role")

	// Fetch the record
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

// Handle List with pagination and search
func HandleUser_rolesList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `json:"page"`
		Limit  *int    `json:"limit"`
		Search *string `json:"search"`
	}

	// If there's data, unmarshal it into req
	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// Set default pagination
	page := 1
	limit := 10
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	// Build query with possible search filter
	q := db.Model(&models.User_roles{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(*req.Search)+"%")
	}

	// Get the total count
	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Fetch the records with pagination and preload related data
	var items []models.User_roles
	if err := q.Preload("User").Preload("Role").Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Prepare and send the response
	response := gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (list)",
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	if len(items) > 0 {
		response["data"] = items
	}

	c.JSON(http.StatusOK, response)
}

// Handle Insert User Role
func HandleUser_rolesInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.User_roles
	// Parse request data
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Check if the user_id exists in the users table
	var userCount int64
	if err := db.Model(&models.Users{}).Where("id = ?", req.UserID).Count(&userCount).Error; err != nil || userCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "user_id does not exist"})
		return
	}

	// Check if the role_id exists in the roles table
	var roleCount int64
	if err := db.Model(&models.Roles{}).Where("id = ?", req.RoleID).Count(&roleCount).Error; err != nil || roleCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "role_id does not exist"})
		return
	}

	// Create new entry
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Preload User and Role
	var result models.User_roles
	if err := db.Preload("User").Preload("Role").First(&result, req.ID).Error; err != nil {
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

// Handle Update User Role
func HandleUser_rolesUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		UserID int64  `json:"user_id"`
		RoleID int64  `json:"role_id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Load m row
	var m models.User_roles
	if err := db.First(&m, "id = ?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Validate FKs only if changing them
	if req.UserID != 0 {
		var cnt int64
		if err := db.Model(&models.Users{}).Where("id = ?", req.UserID).Count(&cnt).Error; err != nil || cnt == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "user_id does not exist"})
			return
		}
	}
	if req.RoleID != 0 {
		var cnt int64
		if err := db.Model(&models.Roles{}).Where("id = ?", req.RoleID).Count(&cnt).Error; err != nil || cnt == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "role_id does not exist"})
			return
		}
	}

	// Uniqueness check for (user_id, role_id) pair excluding this row
	newUserID := m.UserID
	if req.UserID != 0 {
		newUserID = req.UserID
	}
	newRoleID := m.RoleID
	if req.RoleID != 0 {
		newRoleID = req.RoleID
	}
	var dup int64
	if err := db.Model(&models.User_roles{}).
		Where("user_id = ? AND role_id = ? AND id <> ?", newUserID, newRoleID, req.ID).
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
	if req.UserID != 0 {
		updates["user_id"] = req.UserID
	}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}

	// Nothing to change
	if len(updates) == 0 {
		if err := db.Preload("User").Preload("Role").First(&m, "id = ?", req.ID).Error; err != nil {
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
	res := db.Model(&models.User_roles{}).Where("id = ?", req.ID).Updates(updates)
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
	var out models.User_roles
	if err := db.Preload("User").Preload("Role").First(&out, "id = ?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found after update"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (updated)",
		"data":    out,
	})
}

// Handle Delete User Role
func HandleUser_rolesDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
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

	res := db.Delete(&models.User_roles{}, req.ID)
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
