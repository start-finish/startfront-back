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

	// USERS_ROLES API
	"USERS_ROLES_get":    handlers.HandleUsers_rolesGet,
	"USERS_ROLES_list":   handlers.HandleUsers_rolesList,
	"USERS_ROLES_create": handlers.HandleUsers_rolesInsert,
	"USERS_ROLES_update": handlers.HandleUsers_rolesUpdate,
	"USERS_ROLES_delete": handlers.HandleUsers_rolesDelete,

	// USERS API
	"login":        handlers.HandleLogin,
	"sign_up":      handlers.HandleSignUp,
	"USERS_list":   handlers.HandleUsersList,
	"USERS_update": handlers.HandleUsersUpdate,
	"USERS_delete": handlers.HandleUsersDelete,

	// ROLES API
	"ROLES_get":    handlers.HandleRolesGet,
	"ROLES_list":   handlers.HandleRolesList,
	"ROLES_create": handlers.HandleRolesInsert,
	"ROLES_update": handlers.HandleRolesUpdate,
	"ROLES_delete": handlers.HandleRolesDelete,

	// HOME API
	"HOME_get":    handlers.HandleHomeGet,
	"HOME_list":   handlers.HandleHomeList,
	"HOME_create": handlers.HandleHomeInsert,
	"HOME_update": handlers.HandleHomeUpdate,
	"HOME_delete": handlers.HandleHomeDelete,
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
