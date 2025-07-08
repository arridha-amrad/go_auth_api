package routes

import (
	"database/sql"
	"my-go-api/internal/config"
	"my-go-api/internal/controllers/auth"
	"my-go-api/internal/controllers/user"
	"my-go-api/internal/middleware"
	"my-go-api/internal/utils"

	"my-go-api/internal/repositories"
	"my-go-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	db *sql.DB,
	rdb *redis.Client,
	validate *validator.Validate,
	config *config.Config,
) *gin.Engine {

	router := gin.Default()

	utilities := utils.NewUtilities(config.JWtSecretKey, config.AppUri, config.GoogleOAuth2)

	userRepo := repositories.NewUserRepository(db)
	redisRepo := repositories.NewRedisRepository(rdb)

	// services
	redisService := services.NewRedisService(redisRepo)
	jwtService := services.NewJwtService(config.JWtSecretKey, redisService)
	authService := services.NewAuthService(redisService, utilities, jwtService)
	userService := services.NewUserService(userRepo)
	emailService := services.NewEmailService(config.AppUri, utilities)
	passwordService := services.NewPasswordService()

	userController := user.NewUserController(userService)
	authController := auth.NewAuthController(
		passwordService,
		authService,
		userService,
		emailService,
		redisService,
		utilities,
	)

	validationMiddleware := middleware.NewValidationMiddleware(validate)
	authMiddleware := middleware.NewAuthMiddleware(jwtService, userService)

	router.SetTrustedProxies([]string{"127.0.0.1"})

	v1 := router.Group("/api/v1")
	{
		v1.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to V1"})
		})

		SetUserRoutes(UserRoutes{
			route:                v1,
			userController:       userController,
			validationMiddleware: validationMiddleware,
		})

		SetAuthRoutes(AuthRoutesParams{
			route:                v1,
			authController:       authController,
			authMiddleware:       authMiddleware,
			validationMiddleware: validationMiddleware,
		})
	}

	return router
}
