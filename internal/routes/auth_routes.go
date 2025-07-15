package routes

import (
	"my-go-api/internal/controllers/auth"
	"my-go-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRoutesParams struct {
	route                *gin.RouterGroup
	authController       auth.IAuthController
	validationMiddleware middleware.IValidationMiddleware
	authMiddleware       middleware.IAuthMiddleware
}

func SetAuthRoutes(params AuthRoutesParams) {

	authRoutes := params.route.Group("/auth")
	{
		authRoutes.GET("", params.authMiddleware.Handler, params.authController.GetAuth)
		authRoutes.POST("", params.validationMiddleware.Login, params.authController.Login)
		authRoutes.POST("/refresh-token", params.authController.RefreshToken)
		authRoutes.POST("/reset-password", params.validationMiddleware.ResetPassword, params.authController.ResetPassword)
		authRoutes.POST("/forgot-password", params.validationMiddleware.ForgotPassword, params.authController.ForgotPassword)
		authRoutes.POST("/logout", params.authMiddleware.Handler, params.authController.Logout)
		authRoutes.POST("/register", params.validationMiddleware.Register, params.authController.Register)
		authRoutes.POST("/resend-verification", params.validationMiddleware.ResendVerification, params.authController.ResendVerification)
		authRoutes.POST("/verify", params.validationMiddleware.VerifyNewAccount, params.authController.VerifyNewAccount)
	}
}
