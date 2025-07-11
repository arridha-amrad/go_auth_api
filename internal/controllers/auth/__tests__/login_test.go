package auth_test

import (
	"errors"
	"my-go-api/internal/controllers/auth"
	"my-go-api/internal/dto"
	"my-go-api/internal/models"
	"my-go-api/internal/services"
	mockutils "my-go-api/mocks"
	mockservices "my-go-api/mocks/mock_services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mockservices.NewMockIUserService(ctrl)
	mockAuthService := mockservices.NewMockIAuthService(ctrl)
	mockEmailService := mockservices.NewMockIEmailService(ctrl)
	mockPasswordService := mockservices.NewMockIPasswordService(ctrl)
	mockRedisService := mockservices.NewMockIRedisService(ctrl)
	mockUtils := mockutils.NewMockIUtils(ctrl)
	controller := auth.NewAuthController(
		mockPasswordService,
		mockAuthService,
		mockUserService,
		mockEmailService,
		mockRedisService,
		mockUtils,
	)
	gin.SetMode(gin.TestMode)
	// Simulate validated body middleware
	body := dto.Login{
		Identity: "ari@mail.com",
		Password: "password123",
	}
	user := models.User{
		ID:         uuid.New(),
		Username:   "ari00",
		Email:      "ari@mail.com",
		Password:   "hashed-password",
		JwtVersion: "v1",
		IsVerified: true,
		Name:       "Arridha Amrad",
		Provider:   "credentials",
		Role:       "user",
		CreatedAt:  time.Now().String(),
		UpdatedAt:  time.Now().String(),
	}
	authTokens := services.CreateAuthTokensResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
	// Set expectations
	mockUserService.EXPECT().GetUserByIdentity(gomock.Any(), "ari@mail.com").Return(&user, nil)
	mockPasswordService.EXPECT().Verify("hashed-password", "password123").Return(nil)
	mockAuthService.EXPECT().CreateAuthTokens(services.CreateAuthTokenParams{
		UserId:     user.ID,
		JwtVersion: "v1",
	}).Return(authTokens, nil)
	// Setup Gin context
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("validatedBody", body)
	// Run
	controller.Login(c)
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access-token")
	assert.Contains(t, w.Body.String(), "ari@mail.com")
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockservices.NewMockIUserService(ctrl)
	mockAuthService := mockservices.NewMockIAuthService(ctrl)
	mockEmailService := mockservices.NewMockIEmailService(ctrl)
	mockPasswordService := mockservices.NewMockIPasswordService(ctrl)
	mockRedisService := mockservices.NewMockIRedisService(ctrl)
	mockUtils := mockutils.NewMockIUtils(ctrl)

	controller := auth.NewAuthController(
		mockPasswordService,
		mockAuthService,
		mockUserService,
		mockEmailService,
		mockRedisService,
		mockUtils,
	)

	gin.SetMode(gin.TestMode)

	body := dto.Login{
		Identity: "ari@mail.com",
		Password: "password123",
	}

	mockUserService.EXPECT().GetUserByIdentity(gomock.Any(), "ari@mail.com").Return(models.User{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("validatedBody", body)

	controller.Login(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}
