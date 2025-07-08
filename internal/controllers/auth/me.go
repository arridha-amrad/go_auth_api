package auth

import (
	"my-go-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) GetAuth(c *gin.Context) {
	// 1. Extract token payload from context
	value, exist := c.Get("accessTokenPayload")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "validated body not exists"})
		return
	}

	// 2. Type assertion
	tokenPayload, ok := value.(services.VerifyAccessTokenResult)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
		return
	}

	// 3. Fetch user
	user, err := ctrl.userService.GetUserById(c.Request.Context(), tokenPayload.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
