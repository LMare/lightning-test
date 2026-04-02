package testsupport

import (
	"github.com/Lmare/lightning-playground/backend/internal/lightning/port"
	"github.com/stretchr/testify/mock"
)

type MockNodeRepository struct {
	mock.Mock
}

func (m *MockNodeRepository) GetNodesIds() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

var _ port.NodeRepository = (*MockNodeRepository)(nil)
