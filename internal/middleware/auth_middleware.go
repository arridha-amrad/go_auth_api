package middleware

import (
	"my-go-api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	tokenService services.ITokenService
}

type IAuthMiddleware interface {
	Handler(c *gin.Context)
}

func NewAuthMiddleware(tokenService services.ITokenService) IAuthMiddleware {
	return &authMiddleware{
		tokenService: tokenService,
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

	c.Set("accessTokenPayload", payload)

	c.Next()
}
