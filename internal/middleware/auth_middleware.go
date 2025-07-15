package middleware

import (
	"my-go-api/internal/constants"
	"my-go-api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authMiddleware struct {
	userService services.IUserService
	jwtService  services.IJwtService
}

type IAuthMiddleware interface {
	Handler(c *gin.Context)
}

func NewAuthMiddleware(jwtService services.IJwtService, userService services.IUserService) IAuthMiddleware {
	return &authMiddleware{
		userService: userService,
		jwtService:  jwtService,
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

	payload, err := m.jwtService.Verify(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	userId, err := uuid.Parse(payload.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user, err := m.userService.GetUserById(c.Request.Context(), userId)
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

	c.Set(constants.ACCESS_TOKEN_PAYLOAD, payload)

	c.Next()
}
