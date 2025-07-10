package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey    string
	redisService IRedisService
}

type IJwtService interface {
	Verify(tokenString string) (JWTPayload, error)
	Create(JWTPayload) (string, error)
}

func NewJwtService(secretKey string, redisService IRedisService) IJwtService {
	return &jwtService{
		secretKey:    secretKey,
		redisService: redisService,
	}
}

func (s *jwtService) Verify(tokenString string) (JWTPayload, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return JWTPayload{}, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return JWTPayload{}, errors.New("invalid token")
	}

	if _, err := s.redisService.GetAccessToken(claims.JTI); err != nil {
		return JWTPayload{}, err
	}

	return JWTPayload{
		UserId:     claims.UserID,
		Jti:        claims.JTI,
		JwtVersion: claims.JwtVersion,
	}, nil
}

func (s *jwtService) Create(params JWTPayload) (string, error) {
	claims := CustomClaims{
		UserID:     params.UserId,
		JTI:        params.Jti,
		JwtVersion: params.JwtVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// helpers
type CustomClaims struct {
	UserID     string `json:"userId"`
	JTI        string `json:"jti"`
	JwtVersion string `json:"jwtVersion"`
	jwt.RegisteredClaims
}

type JWTPayload struct {
	UserId     string
	Jti        string
	JwtVersion string
}
