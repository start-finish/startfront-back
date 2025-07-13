package pkg

import "github.com/gin-gonic/gin"

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "success"
	}
	c.JSON(200, Response{
		Code:    "0",
		Message: message,
		Data:    data,
	})
}

func ErrorWithStatus(c *gin.Context, message string) {
	c.JSON(400, Response{
		Code:    "1",
		Message: message,
	})
}
