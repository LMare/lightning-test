package repository

import (
	"os"
	"path/filepath"

	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/exception"
)

type NodeRepository interface {
	GetNodesIds() ([]string, error)
}

func NewNodeRepositoryFileSystem(conf *config.Config) NodeRepository {
	return &NodeRepositoryFileSystem{conf: *conf}
}

type NodeRepositoryFileSystem struct {
	conf config.Config
}

// Return the list of node Ids
// Search in the filesystem
func (n *NodeRepositoryFileSystem) GetNodesIds() ([]string, error) {
	nodeIds := []string{}

	// read Folder of Nodes
	entries, err := os.ReadDir(n.conf.NodeStorage)
	if err != nil {
		return nodeIds, exception.NewError("Erreur on reading NodeStorage", err, exception.NewExampleError)
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
