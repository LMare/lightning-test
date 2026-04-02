package testsupport

import (
	"context"

	"github.com/Lmare/lightning-playground/backend/internal/user/domain"
	"github.com/Lmare/lightning-playground/backend/internal/user/port"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]domain.UserModel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.UserModel), args.Error(1)
}

var _ port.UserRepository = (*MockUserRepository)(nil)
