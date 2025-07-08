package auth

import (
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) RefreshToken(c *gin.Context) {
	cookieRefToken, err := c.Cookie(constants.COOKIE_REFRESH_TOKEN)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	data, err := ctrl.tokenService.GetRefreshToken(ctrl.tokenService.HashWithSHA256(cookieRefToken))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println(data)

	user, err := ctrl.userService.GetUserById(c.Request.Context(), data.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	authToken, err := ctrl.tokenService.CreateAuthTokens(services.CreateAuthTokenParams{
		UserId:      data.UserId,
		JwtVersion:  user.JwtVersion,
		OldRefToken: &cookieRefToken,
		OldTokenJti: &data.Jti,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(constants.COOKIE_REFRESH_TOKEN, authToken.RefreshToken, 3600*24*365, "/", "", os.Getenv("GO_ENV") == "production", true)

	c.JSON(http.StatusOK, gin.H{"token": authToken.AccessToken})
}
