package application

import (
	"context"

	"github.com/Lmare/lightning-playground/backend/internal/user/domain"
	"github.com/Lmare/lightning-playground/backend/internal/user/port"
)

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	repo port.UserRepository
}

func (s *UserService) List(ctx context.Context) ([]domain.UserModel, error) {
	return s.repo.FindAll(ctx)
}
