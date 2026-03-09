package testsupport

import (
	"context"

	user "github.com/Lmare/lightning-playground/backend/model/user"
	repository "github.com/Lmare/lightning-playground/backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]user.UserModel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.UserModel), args.Error(1)
}

// Assertion compile-time
var _ repository.UserRepository = (*MockUserRepository)(nil)
