package testsupport

import (
	"context"

	user "github.com/Lmare/lightning-playground/backend/model/user"
	repository "github.com/Lmare/lightning-playground/backend/repository"
)

type FakeUserRepository struct {
	MockFindAll func(ctx context.Context) ([]user.UserModel, error)
}

func (f *FakeUserRepository) FindAll(ctx context.Context) ([]user.UserModel, error) {
	if f.MockFindAll != nil {
		return f.MockFindAll(ctx)
	}
	// Implementation for testing
	return nil, nil
}

// Assertion compile-time
var _ repository.UserRepository = (*FakeUserRepository)(nil)
