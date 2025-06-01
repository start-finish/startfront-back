package handler

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateWidgetBinding(c *gin.Context) {
	var binding domain.WidgetBinding
	if err := c.ShouldBindJSON(&binding); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	if err := usecase.CreateWidgetBinding(binding); err != nil {
		response.Error(c, "Failed to create widget binding")
		return
	}

	response.Success(c, "Widget binding created successfully", nil)
}

func GetWidgetBindings(c *gin.Context) {
	widgetIDStr := c.Param("widget_id")
	widgetID, err := strconv.Atoi(widgetIDStr)
	if err != nil {
		response.Error(c, "Invalid widget ID")
		return
	}

	bindings, err := usecase.GetWidgetBindings(widgetID)
	if err != nil {
		response.Error(c, "Failed to fetch widget bindings")
		return
	}

	response.Success(c, "", gin.H{"data": bindings})
}
