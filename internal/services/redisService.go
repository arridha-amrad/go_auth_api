package services

import (
	"errors"
	"fmt"
	"my-go-api/internal/repositories"
	"time"
)

type redisService struct {
	redisRepository repositories.IRedisRepository
}

type IRedisService interface {
	// access token
	GetAccessToken(jti string) (AccessTokenData, error)
	SaveAccessToken(params AccessTokenData) error
	DeleteAccessToken(jti string) error
	// refresh token
	GetRefreshToken(hashedToken string) (RefreshTokenData, error)
	SaveRefreshToken(params RefreshTokenData) error
	DeleteRefreshToken(hashedToken string) error
	// verification token
	SaveVerificationToken(params VerificationData) error
	DeleteVerificationToken(hashedToken string) error
	GetVerificationToken(hashedToken string) (VerificationData, error)
}

func NewRedisService(redisRepository repositories.IRedisRepository) IRedisService {
	return &redisService{
		redisRepository: redisRepository,
	}
}

func (s *redisService) SaveVerificationToken(params VerificationData) error {
	key := setVerificationKey(params.HashedToken)
	err := s.redisRepository.HSet(key, map[string]any{
		"code":   params.Code,
		"userId": params.UserId,
	}, VerificationTokenTTL)
	return err
}

func (s *redisService) DeleteAccessToken(jti string) error {
	key := setAccessTokenKey(jti)
	err := s.redisRepository.Delete(key)
	return err
}

func (s *redisService) GetRefreshToken(hashedToken string) (RefreshTokenData, error) {
	key := setRefreshTokenKey(hashedToken)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		return RefreshTokenData{}, err
	}

	strUserId, ok := data["userId"]
	if !ok {
		return RefreshTokenData{}, errors.New("userId not found")
	}

	strJti, ok := data["jti"]
	if !ok {
		return RefreshTokenData{}, errors.New("jti not found")
	}

	return RefreshTokenData{
		UserId:      strUserId,
		Jti:         strJti,
		HashedToken: hashedToken,
	}, nil
}

func (s *redisService) DeleteRefreshToken(hashedToken string) error {
	key := setRefreshTokenKey(hashedToken)
	return s.redisRepository.Delete(key)
}

func (s *redisService) DeleteVerificationToken(hashedToken string) error {
	key := setVerificationKey(hashedToken)
	if err := s.redisRepository.Delete(key); err != nil {
		return err
	}
	return nil
}

func (s *redisService) GetVerificationToken(hashedToken string) (VerificationData, error) {
	key := setVerificationKey(hashedToken)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		return VerificationData{}, err
	}

	strCode, ok := data["code"]
	if !ok {
		return VerificationData{}, errors.New("code not found")
	}

	strUserId, ok := data["userId"]
	if !ok {
		return VerificationData{}, errors.New("userId not found")
	}

	return VerificationData{
		Code:   strCode,
		UserId: strUserId,
	}, nil
}

func (s *redisService) SaveRefreshToken(params RefreshTokenData) error {
	key := setRefreshTokenKey(params.HashedToken)
	err := s.redisRepository.HSet(key, map[string]any{
		"userId": params.UserId,
		"jti":    params.Jti,
	}, RefreshTokenTTL)
	return err
}

func (s *redisService) SaveAccessToken(params AccessTokenData) error {
	key := setAccessTokenKey(params.Jti)
	err := s.redisRepository.HSet(key, map[string]any{
		"userId":      params.UserId,
		"accessToken": params.Jti,
	}, AccessTokenTTL)
	return err
}

func (s *redisService) GetAccessToken(jti string) (AccessTokenData, error) {
	key := setAccessTokenKey(jti)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil || len(data) == 0 {
		return AccessTokenData{}, fmt.Errorf("record not found for key : %s", key)
	}
	accessToken, ok := data["accessToken"]
	userId, ok2 := data["userId"]
	if !ok || !ok2 {
		return AccessTokenData{}, errors.New("malformed data")
	}
	return AccessTokenData{
		AccessToken: accessToken,
		UserId:      userId,
		Jti:         jti,
	}, nil
}

// helpers

func setAccessTokenKey(jti string) string {
	return fmt.Sprintf("accessToken:%s", jti)
}

func setRefreshTokenKey(hashedToken string) string {
	return fmt.Sprintf("refreshToken:%s", hashedToken)
}

func setVerificationKey(hashedToken string) string {
	return fmt.Sprintf("account_verification:%s", hashedToken)
}

type RefreshTokenData struct {
	HashedToken string
	UserId      string
	Jti         string
}

type AccessTokenData struct {
	AccessToken string
	UserId      string
	Jti         string
}

type VerificationData struct {
	Code        string
	UserId      string
	HashedToken string
}

var (
	AccessTokenTTL       = 1 * time.Hour
	RefreshTokenTTL      = 24 * 7 * time.Hour
	VerificationTokenTTL = 30 * time.Minute
)
