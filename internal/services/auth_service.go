package services

import (
	"errors"
	"log"
	"my-go-api/internal/utils"

	"github.com/google/uuid"
)

type authService struct {
	redisService IRedisService
	utils        utils.IUtils
	jwtService   IJwtService
}

type IAuthService interface {
	CreateAuthTokens(params CreateAuthTokenParams) (CreateAuthTokensResult, error)
	CreateVerificationToken(userId uuid.UUID) (VerificationTokenData, error)
	VerifyVerificationToken(params VerificationTokenData) (string, error)

	// helpers (not exported)
	GeneratePairToken() (TokenPair, error)
}

func NewAuthService(redisService IRedisService, utils utils.IUtils, jwtService IJwtService) IAuthService {
	return &authService{
		redisService: redisService,
		utils:        utils,
		jwtService:   jwtService,
	}
}

func (s *authService) CreateAuthTokens(params CreateAuthTokenParams) (CreateAuthTokensResult, error) {
	// delete old refresh token record from redis (refresh token behavior)
	if params.OldRefToken != nil {
		if err := s.redisService.DeleteRefreshToken(s.utils.HashWithSHA256(*params.OldRefToken)); err != nil {
			log.Printf("failed to delete refresh token: %s", err.Error())
			return CreateAuthTokensResult{}, err
		}
	}
	// delete old access token record from redis (refresh token behavior)
	if params.OldTokenJti != nil {
		if err := s.redisService.DeleteAccessToken(params.OldTokenJti.String()); err != nil {
			return CreateAuthTokensResult{}, err
		}
	}
	newJti := uuid.New()
	refTokenPair, err := s.GeneratePairToken()
	if err != nil {
		return CreateAuthTokensResult{}, err
	}
	if err := s.redisService.SaveRefreshToken(RefreshTokenData{
		HashedToken: refTokenPair.Hashed,
		UserId:      params.UserId.String(),
		Jti:         newJti.String(),
	}); err != nil {
		log.Println("failed to store refresh token in redis")
		return CreateAuthTokensResult{}, err
	}
	accessToken, err := s.jwtService.Create(JWTPayload{
		UserId:     params.UserId.String(),
		Jti:        newJti.String(),
		JwtVersion: params.JwtVersion,
	})
	if err != nil {
		return CreateAuthTokensResult{}, err
	}
	if err := s.redisService.SaveAccessToken(AccessTokenData{
		AccessToken: accessToken,
		UserId:      params.UserId.String(),
		Jti:         newJti.String(),
	}); err != nil {
		log.Println("failed to store access token in redis")
		return CreateAuthTokensResult{}, err
	}
	return CreateAuthTokensResult{
		RefreshToken: refTokenPair.Raw,
		AccessToken:  accessToken,
	}, nil

}

func (s *authService) CreateVerificationToken(userId uuid.UUID) (VerificationTokenData, error) {
	tokenPair, err := s.GeneratePairToken()
	if err != nil {
		return VerificationTokenData{}, err
	}
	code, err := s.utils.GenerateRandomBytes(4)
	if err != nil {
		return VerificationTokenData{}, err
	}
	if err := s.redisService.SaveVerificationToken(VerificationData{
		Code:        code,
		UserId:      userId.String(),
		HashedToken: tokenPair.Hashed,
	}); err != nil {
		return VerificationTokenData{}, err
	}
	return VerificationTokenData{
		RawToken: tokenPair.Raw,
		Code:     code,
	}, nil
}

func (s *authService) VerifyVerificationToken(params VerificationTokenData) (string, error) {
	data, err := s.redisService.GetVerificationToken(s.utils.HashWithSHA256(params.RawToken))
	if err != nil {
		return "", err
	}
	if data.Code != params.Code {
		return "", errors.New("invalid code")
	}
	return data.UserId, nil
}

// Helpers
func (s *authService) GeneratePairToken() (TokenPair, error) {
	rawToken, err := s.utils.GenerateRandomBytes(32)
	if err != nil {
		return TokenPair{}, errors.New("failure on generating random bytes")
	}

	hashedToken := s.utils.HashWithSHA256(rawToken)
	return TokenPair{
		Raw:    rawToken,
		Hashed: hashedToken,
	}, nil
}

type CreateAuthTokenParams struct {
	UserId      uuid.UUID
	JwtVersion  string
	OldRefToken *string
	OldTokenJti *uuid.UUID
}

type CreateAuthTokensResult struct {
	RefreshToken string
	AccessToken  string
}

type VerificationTokenData struct {
	RawToken string
	Code     string
}

type TokenPair struct {
	Hashed string
	Raw    string
}
