package services

import (
	"my-go-api/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	Raw    string
	Hashed string
}

type VerifyAccessTokenResult struct {
	UserId     uuid.UUID
	Jti        uuid.UUID
	JwtVersion string
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

type CreateAuthTokensResult struct {
	RefreshToken string
	AccessToken  string
}

type CreateAccountVerificationTokenAndCodeResult struct {
	RawToken string
	Code     string
}

type GetRefreshTokenResult struct {
	UserId uuid.UUID
	Jti    uuid.UUID
}

type GetAccessTokenResult struct {
	UserId      uuid.UUID
	AccessToken string
}

type CustomClaims struct {
	UserID     string `json:"userId"`
	JTI        string `json:"jti"`
	JwtVersion string `json:"jwtVersion"`
	jwt.RegisteredClaims
}

type ITokenService interface {
	// Access Token
	VerifyAccessToken(tokenString string) (VerifyAccessTokenResult, error)
	DeleteAccessToken(jti uuid.UUID) error
	// GetAccessToken(jti string) (GetAccessTokenResult, error)

	// Refresh Token
	DeleteRefreshToken(hashedToken string) error
	GetRefreshToken(hashedToken string) (GetRefreshTokenResult, error)

	// Utilities
	HashWithSHA256(randomStr string) string
	GenerateRandomBytes(size int) (string, error)

	// Verification Token
	DeleteVerificationToken(hashedToken string) error
	GetVerificationToken(hashedToken string) (map[string]string, error)

	// Auth Tokens Generation
	CreateAuthTokens(params CreateAuthTokenParams) (CreateAuthTokensResult, error)

	// Account Verification
	CreateAccountVerificationTokenAndCode(userId uuid.UUID) (CreateAccountVerificationTokenAndCodeResult, error)
	VerifyNewAccountTokenAndCode(params CreateAccountVerificationTokenAndCodeResult) (string, error)
}

type tokenService struct {
	redisRepository repositories.IRedisRepository
	secret          string
}
