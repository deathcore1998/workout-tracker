package services

import (
	"context"
	"fmt"

	"github.com/deathcore1998/workout-tracker/internal/models"
	"github.com/deathcore1998/workout-tracker/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password %w", err)
	}

	user := &models.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
