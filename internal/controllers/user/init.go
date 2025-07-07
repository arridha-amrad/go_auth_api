package user

import (
	"my-go-api/internal/services"

	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetUserById(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
}

type userController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) IUserController {
	return &userController{userService: userService}
}
