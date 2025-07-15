package auth

import (
	"fmt"
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"my-go-api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctrl *authController) ForgotPassword(c *gin.Context) {
	value, ok := c.Get(constants.VALIDATED_BODY)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validated body not exists"})
		return
	}

	body, ok := value.(dto.ForgotPassword)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid type for validated body"})
		return
	}

	user, err := ctrl.userService.GetUserByEmail(c.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify your account first"})
		return
	}

	pairToken, err := ctrl.authService.GeneratePairToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	if err := ctrl.redisService.SavePasswordResetToken(services.PasswordResetData{
		HashedToken: pairToken.Hashed,
		UserId:      user.ID.String(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	if err := ctrl.emailService.SendPasswordResetRequest(services.SendPasswordResetParams{
		Name:  user.Username,
		Email: user.Email,
		Token: pairToken.Raw,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("An email has been sent to %s. Please follow the instruction to reset your password", body.Email),
	})

}
