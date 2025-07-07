package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

type IPasswordService interface {
	Verify(hashed, plain string) error
	Hash(password string) (string, error)
}

func NewPasswordService() IPasswordService {
	return &passwordService{}
}

func (s *passwordService) Verify(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func (s *passwordService) Hash(password string) (string, error) {
	// Generate a hashed version of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
