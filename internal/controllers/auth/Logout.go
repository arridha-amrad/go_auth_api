package auth

import (
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) Logout(c *gin.Context) {
	cookieRefToken, err := c.Cookie(constants.COOKIE_REFRESH_TOKEN)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	value, exist := c.Get(constants.ACCESS_TOKEN_PAYLOAD)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "validated body not exists"})
		return
	}

	tokenPayload, ok := value.(services.JWTPayload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	if err := ctrl.redisService.DeleteAccessToken(tokenPayload.Jti); err != nil {
		log.Println(err.Error() + " failed to delete access token")

	}

	hashedToken := ctrl.utils.HashWithSHA256(cookieRefToken)

	err = ctrl.redisService.DeleteRefreshToken(hashedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.SetCookie(constants.COOKIE_REFRESH_TOKEN, "", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "Logout"})
}
