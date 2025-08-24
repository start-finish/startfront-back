package router

import (
	"encoding/json"
	"net/http"

	"startfront/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Request struct {
	MsgID string          `json:"msgId"`
	Data  json.RawMessage `json:"data"`
}

type MsgHandler func(*gorm.DB, json.RawMessage, *gin.Context)

// Keep empty map on init; create-api will inject entries and add handlers import
var MsgHandlers = map[string]MsgHandler{

	// ROLE_PERMISSIONS API
	"ROLE_PERMISSIONS_get":    handlers.HandleRole_permissionsGet,
	"ROLE_PERMISSIONS_list":   handlers.HandleRole_permissionsList,
	"ROLE_PERMISSIONS_create": handlers.HandleRole_permissionsInsert,
	"ROLE_PERMISSIONS_update": handlers.HandleRole_permissionsUpdate,
	"ROLE_PERMISSIONS_delete": handlers.HandleRole_permissionsDelete,

	// PERMISSIONS API
	"PERMISSIONS_get":    handlers.HandlePermissionsGet,
	"PERMISSIONS_list":   handlers.HandlePermissionsList,
	"PERMISSIONS_create": handlers.HandlePermissionsInsert,
	"PERMISSIONS_update": handlers.HandlePermissionsUpdate,
	"PERMISSIONS_delete": handlers.HandlePermissionsDelete,

	// ADMIN_DASHBOARD API
	"ADMIN_DASHBOARD_get":    handlers.HandleAdmin_dashboardGet,
	"ADMIN_DASHBOARD_list":   handlers.HandleAdmin_dashboardList,
	"ADMIN_DASHBOARD_create": handlers.HandleAdmin_dashboardInsert,
	"ADMIN_DASHBOARD_update": handlers.HandleAdmin_dashboardUpdate,
	"ADMIN_DASHBOARD_delete": handlers.HandleAdmin_dashboardDelete,

	// USER_ROLES API
	"USER_ROLES_get":    handlers.HandleUser_rolesGet,
	"USER_ROLES_list":   handlers.HandleUser_rolesList,
	"USER_ROLES_create": handlers.HandleUser_rolesInsert,
	"USER_ROLES_update": handlers.HandleUser_rolesUpdate,
	"USER_ROLES_delete": handlers.HandleUser_rolesDelete,

	// ROLES API
	"ROLES_get":    handlers.HandleRolesGet,
	"ROLES_list":   handlers.HandleRolesList,
	"ROLES_create": handlers.HandleRolesInsert,
	"ROLES_update": handlers.HandleRolesUpdate,
	"ROLES_delete": handlers.HandleRolesDelete,

	// USERS API
	"login":        handlers.HandleLogin,
	"sign_up":      handlers.HandleSignUp,
	"USERS_list":   handlers.HandleUsersList,
	"USERS_update": handlers.HandleUsersUpdate,
	"USERS_delete": handlers.HandleUsersDelete,
}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/api/startProcess", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Invalid JSON",
			})
			return
		}
		handler, exists := MsgHandlers[req.MsgID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Unknown msgId",
			})
			return
		}
		handler(db, req.Data, c)
	})
}
