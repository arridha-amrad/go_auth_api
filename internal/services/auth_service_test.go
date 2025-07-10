package services_test

import (
	mockutils "my-go-api/internal/mocks"
	mockservices "my-go-api/internal/mocks/mock_services"
	"my-go-api/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAuthTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mockservices.NewMockIRedisService(ctrl)
	mockUtils := mockutils.NewMockIUtils(ctrl)
	mockJwt := mockservices.NewMockIJwtService(ctrl)

	service := services.NewAuthService(mockRedis, mockUtils, mockJwt)

	userId := uuid.New()
	rawToken := "raw_refresh_token"
	hashedToken := "hashed_refresh_token"
	accessToken := "access_token"
	jwtVersion := "v1"

	mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
	mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
	mockJwt.EXPECT().Create(gomock.Any()).Return(accessToken, nil)
	mockRedis.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
	mockRedis.EXPECT().SaveAccessToken(gomock.Any()).Return(nil)

	result, err := service.CreateAuthTokens(services.CreateAuthTokenParams{
		UserId:     userId,
		JwtVersion: jwtVersion,
	})

	assert.NoError(t, err)
	assert.Equal(t, rawToken, result.RefreshToken)
	assert.Equal(t, accessToken, result.AccessToken)
}
