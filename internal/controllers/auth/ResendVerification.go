package auth

import (
	"fmt"
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"my-go-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) ResendVerification(c *gin.Context) {
	value, exist := c.Get(constants.VALIDATED_BODY)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validated body not exists"})
		return
	}
	body, ok := value.(dto.ResendVerification)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid type for validated body"})
		return
	}
	user, err := ctrl.userService.GetUserByEmail(c.Request.Context(), body.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not registered"})
		return
	}
	if user.IsVerified {
		c.JSON(http.StatusNotFound, gin.H{"error": "user has been verified"})
		return
	}
	data, err := ctrl.authService.CreateVerificationToken(user.ID)
	if err != nil {
		log.Println("failed to create verification token")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.emailService.SendVerificationEmail(
		services.SendEmailVerificationParams{
			Name:  user.Username,
			Email: user.Email,
			Code:  data.Code,
		},
	); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   data.RawToken,
		"message": fmt.Sprintf("An email has been sent to %s. Please follow the instruction to verify your account.", user.Email)},
	)
}
