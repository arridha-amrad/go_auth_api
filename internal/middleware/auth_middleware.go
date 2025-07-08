package middleware

import (
	"log"
	"my-go-api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	tokenService services.ITokenService
	userService  services.IUserService
}

type IAuthMiddleware interface {
	Handler(c *gin.Context)
}

func NewAuthMiddleware(tokenService services.ITokenService, userService services.IUserService) IAuthMiddleware {
	return &authMiddleware{
		tokenService: tokenService,
		userService:  userService,
	}
}

func (m *authMiddleware) Handler(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		c.Abort()
		return
	}

	tokenStr := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))

	payload, err := m.tokenService.VerifyAccessToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user, err := m.userService.GetUserById(c.Request.Context(), payload.UserId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if payload.JwtVersion != user.JwtVersion {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid jwt version"})
		c.Abort()
		return
	}

	log.Printf("payload Jti : %s", payload.Jti)
	log.Printf("payload JwtVersion : %s", payload.JwtVersion)
	log.Printf("payload UserId : %s", payload.UserId)

	c.Set("accessTokenPayload", payload)

	c.Next()
}
