package filesystem

import (
	"os"
	"path/filepath"

	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/internal/shared/exception"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/port"
)

func NewNodeRepositoryFileSystem(conf *config.Config) port.NodeRepository {
	return &NodeRepositoryFileSystem{conf: *conf}
}

type NodeRepositoryFileSystem struct {
	conf config.Config
}

var _ port.NodeRepository = (*NodeRepositoryFileSystem)(nil)

// Return the list of node Ids
// Search in the filesystem
func (n *NodeRepositoryFileSystem) GetNodesIds() ([]string, error) {
	nodeIds := []string{}

	// read Folder of Nodes
	entries, err := os.ReadDir(n.conf.NodeStorage)
	if err != nil {
		return nodeIds, exception.NewError("Error reading NodeStorage", err, exception.NewExampleError)
	}

	// Parcourir les entrées et filtrer uniquement les dossiers
	for _, entry := range entries {
		if entry.IsDir() {
			if matched, err := filepath.Match(n.conf.NodeNamePattern, entry.Name()); err == nil && matched {
				nodeIds = append(nodeIds, entry.Name())
			}
		}
	}

	return nodeIds, nil
}
