package auth

import (
	"fmt"
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"my-go-api/internal/repositories"
	"my-go-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) Register(c *gin.Context) {
	value, exist := c.Get(constants.VALIDATED_BODY)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	body, ok := value.(dto.Register)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "invalid type for validatedBody"},
		)
		return
	}

	hashedPassword, err := ctrl.passwordService.Hash(body.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to hash password"},
		)
		return
	}

	jwtVersion, err := ctrl.utils.GenerateRandomBytes(8)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"errors": err.Error()})
		return
	}

	user, err := ctrl.userService.Store(c.Request.Context(), repositories.CreateOneParams{
		Name:       body.Name,
		Username:   body.Username,
		Email:      body.Email,
		Password:   hashedPassword,
		JWTVersion: jwtVersion,
	})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"errors": err.Error()})
		return
	}

	data, err := ctrl.authService.CreateVerificationToken(user.ID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"errors": err.Error()})
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

	c.JSON(http.StatusCreated, gin.H{
		"token":   data.RawToken,
		"message": fmt.Sprintf("An email has been sent to %s. Please follow the instruction to verify your account.", user.Email)},
	)
}
