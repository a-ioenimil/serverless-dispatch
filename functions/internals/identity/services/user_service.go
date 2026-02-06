package services

import (
	"context"
	"time"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/domain"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser registers a new user in the system after Cognito verification
func (s *UserService) CreateUser(ctx context.Context, id, email string) error {
	user := domain.User{
		ID:        id,
		Email:     email,
		Role:      domain.RoleMember, // Default role
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
	}

	return s.repo.Save(ctx, user)
}
