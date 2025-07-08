package auth

import (
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"my-go-api/internal/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) Login(c *gin.Context) {
	value, exist := c.Get("validatedBody")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validated body not exists"})
		return
	}
	body, ok := value.(dto.Login)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid type for validated body"})
		return
	}

	user, err := ctrl.userService.GetUserByIdentity(c.Request.Context(), body.Identity)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify your account first"})
		return
	}

	if err := ctrl.passwordService.Verify(user.Password, body.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	authToken, err := ctrl.authService.CreateAuthTokens(services.CreateAuthTokenParams{
		UserId:     user.ID,
		JwtVersion: user.JwtVersion,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(constants.COOKIE_REFRESH_TOKEN, authToken.RefreshToken, 3600*24*365, "/", "", os.Getenv("GO_ENV") == "production", true)

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": authToken.AccessToken,
	})
}
