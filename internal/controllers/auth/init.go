package auth

import (
	"my-go-api/internal/services"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Register(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetAuth(c *gin.Context)
	Login(c *gin.Context)
	VerifyNewAccount(c *gin.Context)
}

type authController struct {
	userService     services.IUserService
	tokenService    services.ITokenService
	emailService    services.IEmailService
	passwordService services.IPasswordService
}

func NewAuthController(
	passwordService services.IPasswordService,
	tokenService services.ITokenService,
	userService services.IUserService,
	emailService services.IEmailService,
) IAuthController {
	return &authController{
		userService:     userService,
		passwordService: passwordService,
		emailService:    emailService,
		tokenService:    tokenService,
	}
}
