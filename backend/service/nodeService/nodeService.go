package nodeService

import (
	"os"
	"path/filepath"

	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/exception"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
)

func NewNodeService(conf config.Config) *NodeService {
	return &NodeService{conf: conf}
}

type NodeService struct {
	conf config.Config
}

// Return the indo to connect to the nodes
func (n *NodeService) ListOfNodes() ([]nodeModel.NodeConfigDescriptor, error) {
	nodes := []nodeModel.NodeConfigDescriptor{}

	ids, err := n.getNodesIds()
	if err != nil {
		return nodes, exception.NewError("Error getNodesIds", err, exception.NewExampleError)
	}

	for _, id := range ids {
		nodes = append(nodes, n.getNodeConfigDescriptor(id))
	}

	return nodes, nil
}

// the instance of NodeConfigDescriptor for the node Id
func (n *NodeService) getNodeConfigDescriptor(id string) nodeModel.NodeConfigDescriptor {
	return nodeModel.NodeConfigDescriptor{Id: id, AuthData: n.GetLndClientAuthData(id)}
}

// the instance of LndClientAuthData for the Id
func (n *NodeService) GetLndClientAuthData(id string) nodeModel.LndClientAuthData {
	certPath := n.conf.NodeStorage + id + "/tls.cert"
	macarronPath := n.conf.NodeStorage + id + "/data/chain/bitcoin/" + n.conf.LndNetwork + "/admin.macaroon"
	url := id + "." + n.conf.LndDomain + ":10009"
	return nodeModel.NewLndClientAuthData(certPath, macarronPath, url)
}

// Return the list of node Ids
// Search in the filesystem
func (n *NodeService) getNodesIds() ([]string, error) {
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
