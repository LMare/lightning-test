package nodeService

import (
	yaml "gopkg.in/yaml.v3"
	"os"
	"fmt"

	config "github.com/Lmare/lightning-test"
	exception "github.com/Lmare/lightning-test/backend/exception"
	lightningService "github.com/Lmare/lightning-test/backend/service/lightningService"
)



type NodeConfigDescriptor struct {
	AuthData	lightningService.LndClientAuthData	`yaml:"data"`
	Id			int									`yaml:"id"`
}

type nodesConfigDescriptor struct {
	Nodes	[]NodeConfigDescriptor	`yaml:nodes`
}

// Read a yaml file to get the resource to connect to the nodes
func ListOfNodes() ([]NodeConfigDescriptor, error) {
	y := config.Load().NodesFileDescriptor

	data, err := os.ReadFile(y)
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lors de la lecture du fichier %s de config des nodes", y), err, exception.NewExampleError)
		return nil, err
	}

	var nodesConf nodesConfigDescriptor
	err = yaml.Unmarshal(data, &nodesConf)
	if err != nil {
		err := exception.NewError("Erreur Unmarshal yaml to NodesConfigDescription", err, exception.NewExampleError)
		return nil, err
	}

	return nodesConf.Nodes, nil
}

// get connection data of a specific node
func GetLndClientAuthData(id int) (lightningService.LndClientAuthData, error) {
	nodes, err := ListOfNodes()
	if err != nil {
		err := exception.NewError("Erreur lors de la récupération des nodes", err, exception.NewExampleError)
		return lightningService.NewLndClientAuthData("", "", ""), err
	}

	for _, node := range nodes {
		if node.Id == id {
			return node.AuthData, nil
		}
	}
	err = exception.NewError(fmt.Sprintf("La node id %d n'existe pas", id), err, exception.NewExampleError)
	return lightningService.NewLndClientAuthData("", "", ""), err
}
