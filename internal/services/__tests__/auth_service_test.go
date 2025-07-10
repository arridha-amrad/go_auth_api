package services_test

import (
	"errors"
	"my-go-api/internal/services"
	mockutils "my-go-api/mocks"
	mockservices "my-go-api/mocks/mock_services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRedis *mockservices.MockIRedisService
	mockUtils *mockutils.MockIUtils
	mockJwt   *mockservices.MockIJwtService
	services  services.IAuthService
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRedis = mockservices.NewMockIRedisService(suite.ctrl)
	suite.mockUtils = mockutils.NewMockIUtils(suite.ctrl)
	suite.mockJwt = mockservices.NewMockIJwtService(suite.ctrl)
	suite.services = services.NewAuthService(suite.mockRedis, suite.mockUtils, suite.mockJwt)
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *AuthServiceTestSuite) TestCreateAuthTokens() {

	suite.Run("It should work for login flow", func() {
		userId := uuid.New()
		rawToken := "raw_refresh_token"
		hashedToken := "hashed_refresh_token"
		accessToken := "access_token"
		jwtVersion := "v1"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockJwt.EXPECT().Create(gomock.Any()).Return(accessToken, nil)
		suite.mockRedis.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
		suite.mockRedis.EXPECT().SaveAccessToken(gomock.Any()).Return(nil)

		result, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:     userId,
			JwtVersion: jwtVersion,
		})

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), rawToken, result.RefreshToken)
		assert.Equal(suite.T(), accessToken, result.AccessToken)
	})

	suite.Run("It should work for refresh token flow", func() {
		userId := uuid.New()
		rawToken := "raw_refresh_token"
		hashedToken := "hashed_refresh_token"
		accessToken := "access_token"
		jwtVersion := "v1"
		oldRefToken := "old ref token"
		oldHashedToken := "old_hashed_refresh_token"
		oldJti := uuid.New()

		suite.mockUtils.EXPECT().HashWithSHA256(oldRefToken).Return(oldHashedToken)
		suite.mockRedis.EXPECT().DeleteRefreshToken(oldHashedToken).Return(nil)
		suite.mockRedis.EXPECT().DeleteAccessToken(oldJti.String()).Return(nil)
		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
		suite.mockJwt.EXPECT().Create(gomock.Any()).Return(accessToken, nil)
		suite.mockRedis.EXPECT().SaveAccessToken(gomock.Any()).Return(nil)

		result, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:      userId,
			JwtVersion:  jwtVersion,
			OldRefToken: &oldRefToken,
			OldTokenJti: &oldJti,
		})

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), rawToken, result.RefreshToken)
		assert.Equal(suite.T(), accessToken, result.AccessToken)
	})

}

func (suite *AuthServiceTestSuite) TestVerificationTokenFlow() {
	suite.Run("Successfully create verification token", func() {
		userId := uuid.New()
		rawToken := "verification_token"
		hashedToken := "hashed_verification_token"
		code := "123456"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockUtils.EXPECT().GenerateRandomBytes(4).Return(code, nil)
		suite.mockRedis.EXPECT().SaveVerificationToken(gomock.Any()).Return(nil)

		result, err := suite.services.CreateVerificationToken(userId)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), rawToken, result.RawToken)
		assert.Equal(suite.T(), code, result.Code)
	})

	suite.Run("Fail to create verification token when Redis fails", func() {
		userId := uuid.New()
		rawToken := "verification_token"
		hashedToken := "hashed_verification_token"
		code := "123456"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockUtils.EXPECT().GenerateRandomBytes(4).Return(code, nil)
		suite.mockRedis.EXPECT().SaveVerificationToken(gomock.Any()).Return(errors.New("redis error"))

		_, err := suite.services.CreateVerificationToken(userId)

		assert.Error(suite.T(), err)
	})
}

func (suite *AuthServiceTestSuite) TestVerifyVerificationToken() {
	suite.Run("Successfully verify token with correct code", func() {
		rawToken := "valid_token"
		hashedToken := "hashed_valid_token"
		code := "123456"
		userId := uuid.New().String()

		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().GetVerificationToken(hashedToken).Return(services.VerificationData{
			Code:   code,
			UserId: userId,
		}, nil)

		returnedUserId, err := suite.services.VerifyVerificationToken(services.VerificationTokenData{
			RawToken: rawToken,
			Code:     code,
		})

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), userId, returnedUserId)
	})

	suite.Run("Fail verification with incorrect code", func() {
		rawToken := "valid_token"
		hashedToken := "hashed_valid_token"
		userId := uuid.New().String()

		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().GetVerificationToken(hashedToken).Return(services.VerificationData{
			Code:   "correct_code",
			UserId: userId,
		}, nil)

		_, err := suite.services.VerifyVerificationToken(services.VerificationTokenData{
			RawToken: rawToken,
			Code:     "wrong_code",
		})

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "invalid code")
	})

	suite.Run("Fail verification when token doesn't exist", func() {
		rawToken := "invalid_token"
		hashedToken := "hashed_invalid_token"

		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().GetVerificationToken(hashedToken).Return(services.VerificationData{}, errors.New("not found"))

		_, err := suite.services.VerifyVerificationToken(services.VerificationTokenData{
			RawToken: rawToken,
			Code:     "123456",
		})

		assert.Error(suite.T(), err)
	})
}

func (suite *AuthServiceTestSuite) TestCreateAuthTokensErrorCases() {
	suite.Run("Fail when unable to generate random bytes", func() {
		userId := uuid.New()
		jwtVersion := "v1"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return("", errors.New("random error"))

		_, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:     userId,
			JwtVersion: jwtVersion,
		})

		assert.Error(suite.T(), err)
	})

	suite.Run("Fail when unable to save refresh token", func() {
		userId := uuid.New()
		rawToken := "raw_refresh_token"
		hashedToken := "hashed_refresh_token"
		jwtVersion := "v1"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().SaveRefreshToken(gomock.Any()).Return(errors.New("redis error"))

		_, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:     userId,
			JwtVersion: jwtVersion,
		})

		assert.Error(suite.T(), err)
	})

	suite.Run("Fail when unable to create JWT", func() {
		userId := uuid.New()
		rawToken := "raw_refresh_token"
		hashedToken := "hashed_refresh_token"
		jwtVersion := "v1"

		suite.mockUtils.EXPECT().GenerateRandomBytes(32).Return(rawToken, nil)
		suite.mockUtils.EXPECT().HashWithSHA256(rawToken).Return(hashedToken)
		suite.mockRedis.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
		suite.mockJwt.EXPECT().Create(gomock.Any()).Return("", errors.New("jwt error"))

		_, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:     userId,
			JwtVersion: jwtVersion,
		})

		assert.Error(suite.T(), err)
	})
}

func (suite *AuthServiceTestSuite) TestRefreshTokenEdgeCases() {
	suite.Run("Fail when unable to delete old refresh token", func() {
		userId := uuid.New()
		jwtVersion := "v1"
		oldRefToken := "old_ref_token"
		oldHashedToken := "old_hashed_token"
		oldJti := uuid.New()

		suite.mockUtils.EXPECT().HashWithSHA256(oldRefToken).Return(oldHashedToken)
		suite.mockRedis.EXPECT().DeleteRefreshToken(oldHashedToken).Return(errors.New("delete error"))

		_, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:      userId,
			JwtVersion:  jwtVersion,
			OldRefToken: &oldRefToken,
			OldTokenJti: &oldJti,
		})

		assert.Error(suite.T(), err)
	})

	suite.Run("Fail when unable to delete old access token", func() {
		userId := uuid.New()
		jwtVersion := "v1"
		oldRefToken := "old_ref_token"
		oldHashedToken := "old_hashed_token"
		oldJti := uuid.New()

		suite.mockUtils.EXPECT().HashWithSHA256(oldRefToken).Return(oldHashedToken)
		suite.mockRedis.EXPECT().DeleteRefreshToken(oldHashedToken).Return(nil)
		suite.mockRedis.EXPECT().DeleteAccessToken(oldJti.String()).Return(errors.New("delete error"))

		_, err := suite.services.CreateAuthTokens(services.CreateAuthTokenParams{
			UserId:      userId,
			JwtVersion:  jwtVersion,
			OldRefToken: &oldRefToken,
			OldTokenJti: &oldJti,
		})

		assert.Error(suite.T(), err)
	})
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
