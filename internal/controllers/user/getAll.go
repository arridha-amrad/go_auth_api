package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *userController) GetAll(c *gin.Context) {
	users, err := ctrl.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
