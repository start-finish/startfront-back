package handler

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func CreateWidgetPreset(c *gin.Context) {
	var preset domain.WidgetPreset
	if err := c.ShouldBindJSON(&preset); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	// err := usecase.CreateWidgetPreset(preset)
	// if err != nil {
	// 	response.Error(c, "Failed to create widget preset")
	// 	return
	// }

	response.Success(c, "Widget preset created successfully", nil)
}

func GetWidgetPresets(c *gin.Context) {
	presets, err := usecase.GetWidgetPresets()
	if err != nil {
		response.Error(c, "Failed to fetch widget presets")
		return
	}

	response.Success(c, "", gin.H{"data": presets})
}

// todo: change get by code ==> get by id