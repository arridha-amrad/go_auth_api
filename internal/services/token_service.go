package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"my-go-api/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	Raw    string
	Hashed string
}

type TokenPayload struct {
	UserId uuid.UUID
	Jti    uuid.UUID
}

type RefreshTokenPayload struct {
	UserId uuid.UUID
	Jti    uuid.UUID
}

type CreateAuthTokenParams struct {
	UserId      uuid.UUID
	JwtVersion  string
	OldRefToken *string
	OldTokenJti *uuid.UUID
}

type AuthTokens struct {
	RefreshToken string
	AccessToken  string
}

type AccountVerificationTokenAndCode struct {
	RawToken string
	Code     string
}

type ITokenService interface {
	VerifyAccessToken(tokenString string) (*TokenPayload, error)
	DeleteAccessToken(jti uuid.UUID) error
	DeleteRefreshToken(hashedToken string) error
	GetRefreshToken(hashedToken string) (map[string]string, error)
	HashWithSHA256(randomStr string) string
	DeleteVerificationToken(hashedToken string) error
	GetVerificationToken(hashedToken string) (map[string]string, error)
	CreateAuthToken(params CreateAuthTokenParams) (AuthTokens, error)
	GenerateRandomBytes(size int) (string, error)
	CreateAccountVerificationTokenAndCode(userId uuid.UUID) (AccountVerificationTokenAndCode, error)
	VerifyNewAccountTokenAndCode(params AccountVerificationTokenAndCode) (string, error)
}

type tokenService struct {
	redisRepository repositories.IRedisRepository
	secret          string
}

func NewTokenService(
	redisRepository repositories.IRedisRepository,
	secret string,
) ITokenService {
	return &tokenService{
		redisRepository: redisRepository,
		secret:          secret,
	}
}

func (s *tokenService) CreateAuthToken(params CreateAuthTokenParams) (AuthTokens, error) {
	if params.OldRefToken != nil {
		if err := s.DeleteRefreshToken(s.HashWithSHA256(*params.OldRefToken)); err != nil {
			log.Printf("failed to delete refresh token: %s", err.Error())
			return AuthTokens{}, err
		}
	}

	if params.OldTokenJti != nil {
		if err := s.DeleteAccessToken(*params.OldTokenJti); err != nil {
			return AuthTokens{}, err
		}
	}

	newJti := uuid.New()

	refTokenPair, err := s.generatePairToken()
	if err != nil {
		return AuthTokens{}, err
	}

	// save refresh token for 7 days
	s.redisRepository.HSet(
		setRefTokenKey(refTokenPair.Hashed),
		map[string]any{
			"userId": params.UserId,
			"jti":    newJti,
		}, time.Duration(time.Hour*24*7),
	)

	accessToken, err := s.generateAccessToken(params.UserId, newJti, params.JwtVersion)
	if err != nil {
		return AuthTokens{}, err
	}

	// save access token for 1 hour
	s.redisRepository.HSet(
		setAccTokenKey(newJti),
		map[string]any{
			"userId":      params.UserId,
			"accessToken": accessToken,
		}, time.Duration(time.Hour),
	)

	return AuthTokens{
		RefreshToken: refTokenPair.Raw,
		AccessToken:  accessToken,
	}, nil

}

func (s *tokenService) VerifyAccessToken(tokenString string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.secret), nil
		})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, ok := (*claims)["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	userId, err := uuid.Parse(getClaimString(claims, "userId"))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	jti, err := uuid.Parse(getClaimString(claims, "jti"))
	if err != nil {
		return nil, fmt.Errorf("invalid JTI: %w", err)
	}

	return &TokenPayload{
		UserId: userId,
		Jti:    jti,
	}, nil
}

func (s *tokenService) DeleteRefreshToken(hashedToken string) error {
	key := setRefTokenKey(hashedToken)
	if err := s.redisRepository.Delete(key); err != nil {
		return err
	}
	return nil
}

func (s *tokenService) GetRefreshToken(hashedToken string) (map[string]string, error) {
	key := setRefTokenKey(hashedToken)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *tokenService) DeleteAccessToken(jti uuid.UUID) error {
	key := setAccTokenKey(jti)
	if err := s.redisRepository.Delete(key); err != nil {
		return err
	}
	return nil
}

func (s *tokenService) HashWithSHA256(randomStr string) string {
	hash := sha256.Sum256([]byte(randomStr))
	return hex.EncodeToString(hash[:])
}

func (s *tokenService) generateAccessToken(userId, jti uuid.UUID, jwtVersion string) (string, error) {
	claims := jwt.MapClaims{
		"userId":     userId.String(),
		"jwtVersion": jwtVersion,
		"jti":        jti.String(),
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *tokenService) generatePairToken() (TokenPair, error) {
	rawToken, err := s.GenerateRandomBytes(32)
	if err != nil {
		return TokenPair{}, errors.New("failure on generating random bytes")
	}
	hashedToken := s.HashWithSHA256(rawToken)
	return TokenPair{
		Raw:    rawToken,
		Hashed: hashedToken,
	}, nil
}

func (s *tokenService) GenerateRandomBytes(size int) (string, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *tokenService) CreateAccountVerificationTokenAndCode(userId uuid.UUID) (AccountVerificationTokenAndCode, error) {
	tokenPair, err := s.generatePairToken()
	if err != nil {
		return AccountVerificationTokenAndCode{}, err
	}

	code, err := s.GenerateRandomBytes(4)
	if err != nil {
		return AccountVerificationTokenAndCode{}, err
	}

	if err := s.redisRepository.HSet(
		setAccountVerificationKey(tokenPair.Hashed),
		map[string]any{
			"code":   code,
			"userId": userId.String(),
		},
		time.Duration(time.Minute*30),
	); err != nil {
		return AccountVerificationTokenAndCode{}, err
	}

	return AccountVerificationTokenAndCode{
		RawToken: tokenPair.Raw,
		Code:     code,
	}, nil
}

func (s *tokenService) VerifyNewAccountTokenAndCode(params AccountVerificationTokenAndCode) (string, error) {
	log.Printf("token : %s", params.RawToken)
	key := setAccountVerificationKey(s.HashWithSHA256(params.RawToken))

	log.Printf("key : %s", key)

	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		log.Println("record not found")
		return "", err
	}

	code, ok := data["code"]
	if !ok {
		return "", errors.New("verification code not found")
	}

	userId, ok := data["userId"]
	if !ok {
		return "", errors.New("user ID not found")
	}

	if code != params.Code {
		return "", errors.New("invalid code")
	}

	return userId, nil
}

func (s *tokenService) DeleteVerificationToken(hashedToken string) error {
	key := setAccountVerificationKey(hashedToken)
	if err := s.redisRepository.Delete(key); err != nil {
		return err
	}
	return nil
}

func (s *tokenService) GetVerificationToken(hashedToken string) (map[string]string, error) {
	key := setAccountVerificationKey(hashedToken)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func setAccTokenKey(jti uuid.UUID) string {
	return fmt.Sprintf("accessToken:%s", jti.String())
}

func setRefTokenKey(hashedToken string) string {
	return fmt.Sprintf("refreshToken:%s", hashedToken)
}

func setAccountVerificationKey(hashedToken string) string {
	return fmt.Sprintf("account_verification:%s", hashedToken)
}

func getClaimString(claims *jwt.MapClaims, key string) string {
	if val, ok := (*claims)[key].(string); ok {
		return val
	}
	return ""
}
