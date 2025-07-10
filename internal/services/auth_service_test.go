package services

import (
	"my-go-api/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestAuthService_CreateAuthTokens(t *testing.T) {
	// Setup
	mockRedis := &MockIRedisService{}
	mockJWT := &MockIJwtService{}
	mockUtils := &utils.MockIUtils{}

	authService := &authService{
		redisService: mockRedis,
		jwtService:   mockJWT,
		utils:        mockUtils,
	}

	// Test cases
	t.Run("successful token creation", func(t *testing.T) {
		userId := uuid.New()
		// jti := uuid.New()
		hashedToken := "hashed-refresh-token"
		// oldRefreshToken := "old-raw-refresh-token"
		rawToken := "raw-refresh-token"
		accessToken := "generated-access-token"

		// Mock expectations
		mockUtils.On("GenerateRandomBytes", 32).Return(rawToken, nil)
		mockUtils.On("HashWithSHA256", rawToken).Return(hashedToken)
		mockJWT.On("Create", mock.Anything).Return(accessToken, nil)
		mockRedis.On("SaveRefreshToken", mock.Anything).Return(nil)
		mockRedis.On("SaveAccessToken", mock.Anything).Return(nil)

		// Call method
		params := CreateAuthTokenParams{
			UserId:     userId,
			JwtVersion: "jwt-version",
		}
		result, err := authService.CreateAuthTokens(params)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, accessToken, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)

		mockUtils.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}
