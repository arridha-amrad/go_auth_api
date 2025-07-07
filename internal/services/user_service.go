package services

import (
	"context"
	"fmt"
	"my-go-api/internal/models"
	"my-go-api/internal/repositories"
	"strings"

	"github.com/google/uuid"
)

type IUserService interface {
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	Store(ctx context.Context, params repositories.CreateOneParams) (*models.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByIdentity(ctx context.Context, identity string) (*models.User, error)
}

type userService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := s.userRepo.UpdateOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *userService) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByIdentity(ctx context.Context, identity string) (*models.User, error) {
	if strings.Contains(identity, "@") {
		return s.GetUserByEmail(ctx, identity)
	} else {
		return s.GetUserByUsername(ctx, identity)
	}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetOne(ctx, repositories.GetOneParams{
		Email: &email,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetOne(ctx, repositories.GetOneParams{Username: &username})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) Store(ctx context.Context, params repositories.CreateOneParams) (*models.User, error) {
	user, err := s.userRepo.CreateOne(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}
