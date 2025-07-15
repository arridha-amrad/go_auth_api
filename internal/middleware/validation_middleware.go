package middleware

import (
	"errors"
	"fmt"
	"log"
	"my-go-api/internal/constants"
	"my-go-api/internal/dto"
	"my-go-api/internal/validation"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type validationMiddleware struct {
	validate *validator.Validate
}

type IValidationMiddleware interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	UpdateUser(c *gin.Context)
	VerifyNewAccount(c *gin.Context)
	ForgotPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
	ResendVerification(c *gin.Context)
}

func NewValidationMiddleware(validate *validator.Validate) IValidationMiddleware {
	return &validationMiddleware{validate: validate}
}

func (m *validationMiddleware) runValidation(c *gin.Context, input any) {
	log.Println(input)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Required params are missing"})
		c.Abort()
		return
	}

	switch v := input.(type) {
	case *dto.ResetPassword:
		if v.Password != v.ConfirmPassword {
			log.Println(v.Password)
			log.Println(v.ConfirmPassword)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "passwords do not match",
			})
			c.Abort()
			return
		}
	}

	if err := m.validate.Struct(input); err != nil {
		var validationErrors validator.ValidationErrors
		var msgErrors = make(map[string]string)
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				message := validation.Messages[e.Tag()]
				if e.Param() != "" {
					msgErrors[strings.ToLower(e.Field())] = fmt.Sprintf(message, e.Param())
				} else {
					msgErrors[strings.ToLower(e.Field())] = message
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrors})
		c.Abort()
		return
	}
}

func (m *validationMiddleware) ForgotPassword(c *gin.Context) {
	var input dto.ForgotPassword
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) ResendVerification(c *gin.Context) {
	var input dto.ResendVerification
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) ResetPassword(c *gin.Context) {
	var input dto.ResetPassword
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) VerifyNewAccount(c *gin.Context) {
	var input dto.VerifyNewAccount
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) Login(c *gin.Context) {
	var input dto.Login
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) Register(c *gin.Context) {
	var input dto.Register
	m.runValidation(c, &input)
	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}

func (m *validationMiddleware) UpdateUser(c *gin.Context) {
	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.Abort()
		return
	}

	valErrors := make(map[string]string)
	if username, exists := input["username"].(string); exists {
		if err := m.validate.Var(username, "required,min=5"); err != nil {
			valErrors["username"] = "a minimum of 5 characters is required"
		}
	}

	if name, exists := input["name"].(string); exists {
		if err := m.validate.Var(name, "required,min=5"); err != nil {
			valErrors["name"] = "a minimum of 5 characters is required"
		}
	}

	if email, exists := input["email"].(string); exists {
		if err := m.validate.Var(email, "required,email"); err != nil {
			valErrors["email"] = "invalid email"
		}
	}

	if password, exists := input["password"].(string); exists {
		if err := m.validate.Var(password, "required,strongPassword"); err != nil {
			valErrors["password"] = "a minimum of 5 characters including an uppercase letter, a lowercase letter, and a number is required"
		}
	}

	if role, exists := input["role"].(string); exists {
		if err := m.validate.Var(role, "required,oneof=user admin"); err != nil {
			valErrors["role"] = "unrecognized role"
		}
	}

	if len(valErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": valErrors})
		c.Abort()
		return
	}

	c.Set(constants.VALIDATED_BODY, input)
	c.Next()
}
