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

	// WIDGET_INSTANCES API
	"WIDGET_INSTANCES_get":    handlers.HandleWidget_instancesGet,
	"WIDGET_INSTANCES_list":   handlers.HandleWidget_instancesList,
	"WIDGET_INSTANCES_create": handlers.HandleWidget_instancesInsert,
	"WIDGET_INSTANCES_update": handlers.HandleWidget_instancesUpdate,
	"WIDGET_INSTANCES_delete": handlers.HandleWidget_instancesDelete,

	// WIDGETS API
	"WIDGETS_get":    handlers.HandleWidgetsGet,
	"WIDGETS_list":   handlers.HandleWidgetsList,
	"WIDGETS_create": handlers.HandleWidgetsInsert,
	"WIDGETS_update": handlers.HandleWidgetsUpdate,
	"WIDGETS_delete": handlers.HandleWidgetsDelete,

	// NAVIGATION_ITEMS API
	"NAVIGATION_ITEMS_get":    handlers.HandleNavigation_itemsGet,
	"NAVIGATION_ITEMS_list":   handlers.HandleNavigation_itemsList,
	"NAVIGATION_ITEMS_create": handlers.HandleNavigation_itemsInsert,
	"NAVIGATION_ITEMS_update": handlers.HandleNavigation_itemsUpdate,
	"NAVIGATION_ITEMS_delete": handlers.HandleNavigation_itemsDelete,

	// NAVIGATION_MENUS API
	"NAVIGATION_MENUS_get":    handlers.HandleNavigation_menusGet,
	"NAVIGATION_MENUS_list":   handlers.HandleNavigation_menusList,
	"NAVIGATION_MENUS_create": handlers.HandleNavigation_menusInsert,
	"NAVIGATION_MENUS_update": handlers.HandleNavigation_menusUpdate,
	"NAVIGATION_MENUS_delete": handlers.HandleNavigation_menusDelete,

	// SCREENS API
	"SCREENS_get":    handlers.HandleScreensGet,
	"SCREENS_list":   handlers.HandleScreensList,
	"SCREENS_create": handlers.HandleScreensInsert,
	"SCREENS_update": handlers.HandleScreensUpdate,
	"SCREENS_delete": handlers.HandleScreensDelete,

	// CLIENT_USERS API
	"CLIENT_USERS_get":    handlers.HandleClient_usersGet,
	"CLIENT_USERS_list":   handlers.HandleClient_usersList,
	"CLIENT_USERS_create": handlers.HandleClient_usersInsert,
	"CLIENT_USERS_update": handlers.HandleClient_usersUpdate,
	"CLIENT_USERS_delete": handlers.HandleClient_usersDelete,

	// PROJECTS API
	"PROJECTS_get":    handlers.HandleProjectsGet,
	"PROJECTS_list":   handlers.HandleProjectsList,
	"PROJECTS_create": handlers.HandleProjectsInsert,
	"PROJECTS_update": handlers.HandleProjectsUpdate,
	"PROJECTS_delete": handlers.HandleProjectsDelete,

	// CLIENTS API
	"CLIENTS_get":    handlers.HandleClientsGet,
	"CLIENTS_list":   handlers.HandleClientsList,
	"CLIENTS_create": handlers.HandleClientsInsert,
	"CLIENTS_update": handlers.HandleClientsUpdate,
	"CLIENTS_delete": handlers.HandleClientsDelete,

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
