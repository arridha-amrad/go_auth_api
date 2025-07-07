package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *userController) GetUserById(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id"})
		return
	}

	user, err := ctrl.userService.GetUserById(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
