package testsupport

import (
	repository "github.com/Lmare/lightning-playground/backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockNodeRepository struct {
	mock.Mock
}

func (m *MockNodeRepository) GetNodesIds() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// Assertion compile-time
var _ repository.NodeRepository = (*MockNodeRepository)(nil)
