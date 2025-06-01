package delivery

import (
	"database/sql"
	"startfront-backend/internal/handler"
	"startfront-backend/internal/repository"
	"startfront-backend/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) {
	r := gin.Default()

	// Setup Auth components
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	// Register all API routes under one group
	api := r.Group("/api")
	{
		// Auth
		api.POST("/login", authHandler.Login)

		// WebSocket
		api.GET("/ws", handler.WebSocketHandlerGin)
		go func() {
			time.Sleep(2 * time.Second)
			handler.SendToClients("🔄 Server restarted, welcome back!")
		}()

		// Users
		api.POST("/users", handler.CreateUser)
		api.GET("/users/:id", handler.GetUser)
		api.PUT("/users/:id", handler.UpdateUser)
		api.DELETE("/users/:id", handler.DeleteUser)

		// Applications
		api.POST("/applications", handler.CreateApplication)
		api.GET("/applications", handler.ListApplications)
		api.GET("/application/:id", handler.GetApplication)
		api.PUT("/application/:id", handler.UpdateApplication)
		api.DELETE("/application/:id", handler.DeleteApplication)

		// Screens
		api.POST("/screens", handler.CreateScreen)
		api.GET("/screens", handler.ListScreens)
		api.GET("/screens/:id", handler.GetScreensById)
		api.PUT("/screens/:id", handler.UpdateScreen)
		api.DELETE("/screens/:id", handler.DeleteScreen)

		// Widgets
		api.POST("/widgets", handler.CreateWidget)
		api.GET("/widgets/:screen_id", handler.GetWidgetsByScreenID)
		api.PUT("/widgets/:id", handler.UpdateWidget)
		api.DELETE("/widgets/:id", handler.DeleteWidget)

		// Widget Presets
		api.POST("/widget-presets", handler.CreateWidgetPreset)
		api.GET("/widget-presets", handler.GetWidgetPresets)

		// App Connections
		api.POST("/app-connections", handler.CreateAppConnection)
		api.GET("/app-connections/:application_id", handler.GetAppConnectionsByAppID)

		// Widget Bindings
		api.POST("/widget-bindings", handler.CreateWidgetBinding)
		api.GET("/widget-bindings/:widget_id", handler.GetWidgetBindings)

		// Application Collaborators
		api.POST("/application-collaborators", handler.CreateApplicationCollaborator)
		api.GET("/application-collaborators/:application_id", handler.GetApplicationCollaborators)
	}

	r.Run(":8080")
}
