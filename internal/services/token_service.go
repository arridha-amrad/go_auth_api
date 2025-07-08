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

func NewTokenService(redisRepository repositories.IRedisRepository, secret string) ITokenService {
	return &tokenService{
		redisRepository: redisRepository,
		secret:          secret,
	}
}

func (s *tokenService) CreateAuthTokens(params CreateAuthTokenParams) (CreateAuthTokensResult, error) {
	if params.OldRefToken != nil {
		if err := s.DeleteRefreshToken(s.HashWithSHA256(*params.OldRefToken)); err != nil {
			log.Printf("failed to delete refresh token: %s", err.Error())
			return CreateAuthTokensResult{}, err
		}
	}

	if params.OldTokenJti != nil {
		if err := s.DeleteAccessToken(*params.OldTokenJti); err != nil {
			return CreateAuthTokensResult{}, err
		}
	}

	newJti := uuid.New()

	refTokenPair, err := s.generatePairToken()
	if err != nil {
		return CreateAuthTokensResult{}, err
	}

	// save refresh token for 7 days
	s.redisRepository.HSet(
		setRefTokenKey(refTokenPair.Hashed),
		map[string]any{
			"userId": params.UserId.String(),
			"jti":    newJti.String(),
		}, time.Duration(time.Hour*24*7),
	)

	accessToken, err := s.generateAccessToken(params.UserId, newJti, params.JwtVersion)
	if err != nil {
		return CreateAuthTokensResult{}, err
	}

	// save access token for 1 hour
	s.redisRepository.HSet(
		setAccTokenKey(newJti),
		map[string]any{
			"userId":      params.UserId.String(),
			"accessToken": accessToken,
		}, time.Duration(time.Hour),
	)

	return CreateAuthTokensResult{
		RefreshToken: refTokenPair.Raw,
		AccessToken:  accessToken,
	}, nil

}

func (s *tokenService) VerifyAccessToken(tokenString string) (VerifyAccessTokenResult, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return VerifyAccessTokenResult{}, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return VerifyAccessTokenResult{}, errors.New("invalid token")
	}

	userId, err := uuid.Parse(claims.UserID)
	if err != nil {
		return VerifyAccessTokenResult{}, fmt.Errorf("invalid user ID: %w", err)
	}

	jti, err := uuid.Parse(claims.JTI)
	if err != nil {
		return VerifyAccessTokenResult{}, fmt.Errorf("invalid JTI: %w", err)
	}

	data, err := s.redisRepository.HGetAll(setAccTokenKey(jti))
	if err != nil || len(data) == 0 {
		return VerifyAccessTokenResult{}, errors.New("access token not found or revoked")
	}

	return VerifyAccessTokenResult{
		UserId:     userId,
		Jti:        jti,
		JwtVersion: claims.JwtVersion,
	}, nil
}

func (s *tokenService) DeleteRefreshToken(hashedToken string) error {
	key := setRefTokenKey(hashedToken)
	return s.redisRepository.Delete(key)
}

func (s *tokenService) GetRefreshToken(hashedToken string) (GetRefreshTokenResult, error) {
	key := setRefTokenKey(hashedToken)
	data, err := s.redisRepository.HGetAll(key)
	if err != nil {
		return GetRefreshTokenResult{}, err
	}

	strUserId, ok := data["userId"]
	strJti, ok2 := data["jti"]
	if !ok || !ok2 {
		return GetRefreshTokenResult{}, err
	}

	userId, err := uuid.Parse(strUserId)
	if err != nil {
		return GetRefreshTokenResult{}, err
	}

	jti, err := uuid.Parse(strJti)
	if err != nil {
		return GetRefreshTokenResult{}, err
	}

	return GetRefreshTokenResult{
		UserId: userId,
		Jti:    jti,
	}, nil
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
	claims := CustomClaims{
		UserID:     userId.String(),
		JTI:        jti.String(),
		JwtVersion: jwtVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
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

func (s *tokenService) CreateAccountVerificationTokenAndCode(userId uuid.UUID) (CreateAccountVerificationTokenAndCodeResult, error) {
	tokenPair, err := s.generatePairToken()
	if err != nil {
		return CreateAccountVerificationTokenAndCodeResult{}, err
	}

	code, err := s.GenerateRandomBytes(4)
	if err != nil {
		return CreateAccountVerificationTokenAndCodeResult{}, err
	}

	if err := s.redisRepository.HSet(
		setAccountVerificationKey(tokenPair.Hashed),
		map[string]any{
			"code":   code,
			"userId": userId.String(),
		},
		time.Duration(time.Minute*30),
	); err != nil {
		return CreateAccountVerificationTokenAndCodeResult{}, err
	}

	return CreateAccountVerificationTokenAndCodeResult{
		RawToken: tokenPair.Raw,
		Code:     code,
	}, nil
}

func (s *tokenService) VerifyNewAccountTokenAndCode(params CreateAccountVerificationTokenAndCodeResult) (string, error) {
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
