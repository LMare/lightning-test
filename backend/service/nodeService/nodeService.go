package nodeService

import (
	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/exception"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
	repository "github.com/Lmare/lightning-playground/backend/repository"
)

func NewNodeService(conf config.Config, repo repository.NodeRepository) *NodeService {
	return &NodeService{conf: conf, repo: repo}
}

type NodeService struct {
	conf config.Config
	repo repository.NodeRepository
}

// Return the indo to connect to the nodes
func (n *NodeService) ListOfNodes() ([]nodeModel.NodeConfigDescriptor, error) {
	nodes := []nodeModel.NodeConfigDescriptor{}

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
