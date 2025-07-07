package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *userController) Update(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id"})
		return
	}

	value, exist := c.Get("validatedBody")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	existingUser, err := ctrl.userService.GetUserById(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
	}

	if v, ok := value.(map[string]any); ok {
		if username, exists := v["username"].(string); exists {
			existingUser.Username = username
		}
		if name, exists := v["name"].(string); exists {
			existingUser.Name = name
		}
		if email, exists := v["email"].(string); exists {
			existingUser.Email = email
		}
		if password, exists := v["password"].(string); exists {
			existingUser.Password = password
		}
		if role, exists := v["role"].(string); exists {
			existingUser.Role = role
		}
	}

	ctrl.userService.UpdateUser(c.Request.Context(), existingUser)

	c.JSON(http.StatusOK, gin.H{"user": existingUser})
}
