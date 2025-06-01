package handler

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err := usecase.CreateUser(user)
	if err != nil {
		// Check if the error is due to a duplicate email or name
		if err.Error() == "duplicate_found" {
			response.Error(c, "Duplicate email or name found")
		} else {
			// General error message
			response.Error(c, "Failed to create user")
		}
		return
	}

	// Send WebSocket message
	SendToClients("New user created")

	response.Success(c, "User created successfully", nil)
}

// GetUser fetches user details by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := usecase.GetUser(id)
	if err != nil {
		response.Error(c, "Failed to fetch user")
		return
	}

	response.Success(c, "", gin.H{"user": user})
}

// UpdateUser updates user details
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user domain.User

	_, err := usecase.GetUser(id)
	if err != nil {
		response.Error(c, "User not found")
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err = usecase.UpdateUser(id, user)
	if err != nil {
		response.Error(c, "Failed to update user")
		return
	}

	// Send WebSocket message
	SendToClients("User updated")

	response.Success(c, "User updated successfully", nil)
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := usecase.GetUser(id)
	if err != nil {
		response.Error(c, "User not found")
		return
	}

	err = usecase.DeleteUser(id)
	if err != nil {
		response.Error(c, "Failed to delete user")
		return
	}

	// Send WebSocket message
	SendToClients("User deleted")

	response.Success(c, "User deleted successfully", nil)
}
