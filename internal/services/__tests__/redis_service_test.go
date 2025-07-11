package services_test

import (
	"my-go-api/internal/services"
	mockrepositories "my-go-api/mocks/mock_repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RedisServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockRedisRepo *mockrepositories.MockIRedisRepository
	redisService  services.IRedisService
	sampleAccess  services.AccessTokenData
	sampleRefresh services.RefreshTokenData
	sampleVerify  services.VerificationData
}

func (suite *RedisServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRedisRepo = mockrepositories.NewMockIRedisRepository(suite.ctrl)
	suite.redisService = services.NewRedisService(suite.mockRedisRepo)

	suite.sampleAccess = services.AccessTokenData{
		UserId:      "user123",
		Jti:         "jti-abc",
		AccessToken: "token-xyz",
	}
	suite.sampleRefresh = services.RefreshTokenData{
		UserId:      "user123",
		Jti:         "jti-abc",
		HashedToken: "hashed-refresh",
	}
	suite.sampleVerify = services.VerificationData{
		Code:        "123456",
		UserId:      "user123",
		HashedToken: "hashed-verification",
	}
}

func (suite *RedisServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *RedisServiceTestSuite) TestSaveAccessToken() {
	key := "accessToken:jti-abc"
	suite.mockRedisRepo.EXPECT().HSet(key, gomock.Any(), services.AccessTokenTTL).Return(nil)
	err := suite.redisService.SaveAccessToken(suite.sampleAccess)
	assert.NoError(suite.T(), err)
}

func (suite *RedisServiceTestSuite) TestGetAccessToken_Success() {
	key := "accessToken:jti-abc"
	mockData := map[string]string{
		"userId":      "user123",
		"accessToken": "token-xyz",
	}
	suite.mockRedisRepo.EXPECT().HGetAll(key).Return(mockData, nil)

	data, err := suite.redisService.GetAccessToken("jti-abc")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.sampleAccess, data)
}

func (suite *RedisServiceTestSuite) TestGetAccessToken_RecordNotFound() {
	key := "accessToken:jti-abc"
	suite.mockRedisRepo.EXPECT().HGetAll(key).Return(map[string]string{}, nil)

	_, err := suite.redisService.GetAccessToken("jti-abc")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "record not found")
}

func (suite *RedisServiceTestSuite) TestSaveRefreshToken() {
	key := "refreshToken:hashed-refresh"
	suite.mockRedisRepo.EXPECT().HSet(key, gomock.Any(), services.RefreshTokenTTL).Return(nil)
	err := suite.redisService.SaveRefreshToken(suite.sampleRefresh)
	assert.NoError(suite.T(), err)
}

func (suite *RedisServiceTestSuite) TestGetRefreshToken_Success() {
	key := "refreshToken:hashed-refresh"
	mockData := map[string]string{
		"userId": "user123",
		"jti":    "jti-abc",
	}
	suite.mockRedisRepo.EXPECT().HGetAll(key).Return(mockData, nil)

	data, err := suite.redisService.GetRefreshToken("hashed-refresh")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.sampleRefresh, data)
}

func (suite *RedisServiceTestSuite) TestGetRefreshToken_FieldMissing() {
	key := "refreshToken:hashed-refresh"
	suite.mockRedisRepo.EXPECT().HGetAll(key).Return(map[string]string{
		"userId": "user123",
	}, nil)

	_, err := suite.redisService.GetRefreshToken("hashed-refresh")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "jti not found", err.Error())
}

func (suite *RedisServiceTestSuite) TestSaveVerificationToken() {
	key := "account_verification:hashed-verification"
	suite.mockRedisRepo.EXPECT().HSet(key, gomock.Any(), services.VerificationTokenTTL).Return(nil)
	err := suite.redisService.SaveVerificationToken(suite.sampleVerify)
	assert.NoError(suite.T(), err)
}

func (suite *RedisServiceTestSuite) TestGetVerificationToken_Success() {
	key := "account_verification:hashed-verification"
	mockData := map[string]string{
		"userId": "user123",
		"code":   "123456",
	}
	suite.mockRedisRepo.EXPECT().HGetAll(key).Return(mockData, nil)

	data, err := suite.redisService.GetVerificationToken("hashed-verification")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.sampleVerify.Code, data.Code)
	assert.Equal(suite.T(), suite.sampleVerify.UserId, data.UserId)
}

func (suite *RedisServiceTestSuite) TestDeleteTokens() {
	accessKey := "accessToken:jti-abc"
	refreshKey := "refreshToken:hashed-refresh"
	verifyKey := "account_verification:hashed-verification"

	suite.mockRedisRepo.EXPECT().Delete(accessKey).Return(nil)
	suite.mockRedisRepo.EXPECT().Delete(refreshKey).Return(nil)
	suite.mockRedisRepo.EXPECT().Delete(verifyKey).Return(nil)

	assert.NoError(suite.T(), suite.redisService.DeleteAccessToken("jti-abc"))
	assert.NoError(suite.T(), suite.redisService.DeleteRefreshToken("hashed-refresh"))
	assert.NoError(suite.T(), suite.redisService.DeleteVerificationToken("hashed-verification"))
}

func TestRedisServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RedisServiceTestSuite))
}
