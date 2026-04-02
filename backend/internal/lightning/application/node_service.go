package application

import (
	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/internal/shared/exception"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/domain"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/port"
)

func NewNodeService(conf config.Config, repo port.NodeRepository) *NodeService {
	return &NodeService{conf: conf, repo: repo}
}

type NodeService struct {
	conf config.Config
	repo port.NodeRepository
}

// Return the info to connect to the nodes
func (n *NodeService) ListOfNodes() ([]domain.NodeConfigDescriptor, error) {
	nodes := []domain.NodeConfigDescriptor{}

	ids, err := n.repo.GetNodesIds()
	if err != nil {
		return nodes, exception.NewError("Error getNodesIds", err, exception.NewExampleError)
	}

	for _, id := range ids {
		nodes = append(nodes, n.getNodeConfigDescriptor(id))
	}

	return nodes, nil
}

// the instance of NodeConfigDescriptor for the node Id
func (n *NodeService) getNodeConfigDescriptor(id string) domain.NodeConfigDescriptor {
	return domain.NodeConfigDescriptor{Id: id, AuthData: n.GetLndClientAuthData(id)}
}

// the instance of LndClientAuthData for the Id
func (n *NodeService) GetLndClientAuthData(id string) domain.LndClientAuthData {
	certPath := n.conf.NodeStorage + id + "/tls.cert"
	macarronPath := n.conf.NodeStorage + id + "/data/chain/bitcoin/" + n.conf.LndNetwork + "/admin.macaroon"
	url := id + "." + n.conf.LndDomain + ":10009"
	return domain.NewLndClientAuthData(certPath, macarronPath, url)
}
