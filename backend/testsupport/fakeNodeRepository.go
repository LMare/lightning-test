package testsupport

import (
	repository "github.com/Lmare/lightning-playground/backend/repository"
)

type FakeNodeRepository struct {
	MockGetNodesIds func() ([]string, error)
}

func (f *FakeNodeRepository) GetNodesIds() ([]string, error) {
	if f.MockGetNodesIds != nil {
		return f.MockGetNodesIds()
	}
	// Implementation for testing
	return nil, nil
}

// Assertion compile-time
var _ repository.NodeRepository = (*FakeNodeRepository)(nil)
