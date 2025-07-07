package routes

import (
	"my-go-api/internal/controllers/user"
	"my-go-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	route                *gin.RouterGroup
	userController       user.IUserController
	validationMiddleware middleware.IValidationMiddleware
}

func SetUserRoutes(params UserRoutes) {
	v1Users := params.route.Group("/users")
	{
		v1Users.GET("", params.userController.GetAll)
		v1Users.GET("/:id", params.userController.GetUserById)
		v1Users.PUT("/:id", params.validationMiddleware.UpdateUser, params.userController.Update)
	}
}
