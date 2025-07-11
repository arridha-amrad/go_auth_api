package services_test

import (
	"errors"
	"my-go-api/internal/services"
	mockservices "my-go-api/mocks/mock_services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type JwtServiceTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockRedis   *mockservices.MockIRedisService
	jwtService  services.IJwtService
	secretKey   string
	testPayload services.JWTPayload
}

func (suite *JwtServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRedis = mockservices.NewMockIRedisService(suite.ctrl)
	suite.secretKey = "test-secret"
	suite.jwtService = services.NewJwtService(suite.secretKey, suite.mockRedis)

	suite.testPayload = services.JWTPayload{
		UserId:     "user123",
		Jti:        "jti-abc",
		JwtVersion: "v1",
	}
}

func (suite *JwtServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *JwtServiceTestSuite) TestCreateAndVerifyToken_Success() {
	token, err := suite.jwtService.Create(suite.testPayload)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// Setup redis mock to simulate token exists
	suite.mockRedis.EXPECT().GetAccessToken(suite.testPayload.Jti).Return(services.AccessTokenData{
		AccessToken: token,
		UserId:      suite.testPayload.UserId,
		Jti:         suite.testPayload.Jti,
	}, nil)

	// Verify token
	claims, err := suite.jwtService.Verify(token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testPayload.UserId, claims.UserId)
	assert.Equal(suite.T(), suite.testPayload.Jti, claims.Jti)
	assert.Equal(suite.T(), suite.testPayload.JwtVersion, claims.JwtVersion)
}

func (suite *JwtServiceTestSuite) TestVerify_InvalidToken() {
	invalidToken := "invalid.token.value"

	_, err := suite.jwtService.Verify(invalidToken)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "token parsing failed")
}

func (suite *JwtServiceTestSuite) TestVerify_TokenNotInRedis() {
	// Create valid token
	token, err := suite.jwtService.Create(suite.testPayload)
	assert.NoError(suite.T(), err)

	// Redis returns error (simulating token not found or revoked)
	suite.mockRedis.EXPECT().GetAccessToken(suite.testPayload.Jti).Return(services.AccessTokenData{}, errors.New("not found"))

	_, err = suite.jwtService.Verify(token)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not found", err.Error())
}

func TestJwtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceTestSuite))
}
