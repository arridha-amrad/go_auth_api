package auth

import (
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *authController) ResetPassword(c *gin.Context) {
	value, exist := c.Get(constants.VALIDATED_BODY)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validated body not exists"})
		return
	}
	body, ok := value.(dto.ResetPassword)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid type for validated body"})
		return
	}

	data, err := ctrl.redisService.GetPasswordResetToken(ctrl.utils.HashWithSHA256(body.Token))
	if err != nil {
		log.Println("token not found in redis")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userId, err := uuid.Parse(data.UserId)
	if err != nil {
		log.Println("failed to parse to uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.GetUserById(c.Request.Context(), userId)
	if err != nil {
		log.Println("user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	nv, err := ctrl.utils.GenerateRandomBytes(8)
	if err != nil {
		log.Println("failed to generate new jwt version")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	newPassword, err := ctrl.passwordService.Hash(body.Password)
	if err != nil {
		log.Println("failed to hash the password")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.JwtVersion = nv
	user.Password = newPassword

	if _, err := ctrl.userService.UpdateUser(c.Request.Context(), user); err != nil {
		log.Println("failed to update user jwt_version")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.redisService.DeletePasswordResetToken(ctrl.utils.HashWithSHA256(body.Token)); err != nil {
		log.Println("failed to delete pwd reset token from redis")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset password is successful"})

}
