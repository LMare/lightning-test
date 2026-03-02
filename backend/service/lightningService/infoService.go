package lightningService

import (
	"context"
	"fmt"

	exception "github.com/Lmare/lightning-playground/backend/exception"
	lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
)

func NewLightningInfoService(factory GrpcClientFactory) *LightningInfoService {
	return &LightningInfoService{factory: factory}
}

type LightningInfoService struct {
	factory GrpcClientFactory
}

// Get the active list of node
func (s *LightningInfoService) GetListOfNode(descriptors []nodeModel.NodeConfigDescriptor) []NodeBasicInfo {
	l := []NodeBasicInfo{}

	for _, descriptor := range descriptors {
		basicInfo, err := s.getBasicInfo(descriptor)
		if err != nil {
			fmt.Println("[WARN] ", err)
			continue
		}
		l = append(l, basicInfo)
	}
	return l
}

// get the basical info of a node
func (s *LightningInfoService) getBasicInfo(descriptor nodeModel.NodeConfigDescriptor) (NodeBasicInfo, error) {
	factory, err := s.factory.newClient(descriptor.AuthData)
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Unable to open dial with Node[%s]", descriptor.Id), err, exception.NewExampleError)
		return NodeBasicInfo{}, err
	}
	defer factory.Close()

	resp, err := factory.GetLightningClient().GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Unable to getInfo of Node[%s]", descriptor.Id), err, exception.NewExampleError)
		return NodeBasicInfo{}, err
	}

	return NodeBasicInfo{Id: descriptor.Id, Alias: resp.GetAlias(), Color: resp.GetColor()}, nil
}

// get the uri of the lnd
func (s *LightningInfoService) GetFirstUri(dataClient nodeModel.LndClientAuthData) (string, error) {
	factory, err := s.factory.newClient(dataClient)
	if err != nil {
		err := exception.NewError("Unable to open dial", err, exception.NewExampleError)
		return "", err
	}
	defer factory.Close()

	resp, err := factory.GetLightningClient().GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		err := exception.NewError("Unable to getInfo of Node", err, exception.NewExampleError)
		return "", err
	}
	// extract the uri
	uris := resp.GetUris()
	uri := ""
	if len(uris) > 0 {
		uri = uris[0]
	}

	return uri, nil
}

// return node Information
func (s *LightningInfoService) GetUsefullInfo(dataClient nodeModel.LndClientAuthData) (*InfoLndNode, error) {
	factory, err := s.factory.newClient(dataClient)
	if err != nil {
		err := exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
		return nil, err
	}
	defer factory.Close()

	resp, err := factory.GetLightningClient().GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		err := exception.NewError("Lightning Node respond an error", err, exception.NewExampleError)
		return nil, err
	}

	return &InfoLndNode{
		resp.GetAlias(),
		resp.GetColor(),
		resp.GetNumPendingChannels(),
		resp.GetNumActiveChannels(),
		resp.GetNumInactiveChannels(),
		resp.GetNumPeers(),
		resp.GetBlockHeight(),
		resp.GetChains()[0].GetNetwork(),
		resp.GetUris(),
		resp.GetSyncedToChain(),
		resp.GetSyncedToGraph(),
	}, nil
}

// change the color and the alias of the node
func (s *LightningInfoService) UpdateAliasAndColor(dataClient nodeModel.LndClientAuthData, alias string, color string) error {
	factory, err := s.factory.newClient(dataClient)
	if err != nil {
		err := exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
		return err
	}
	defer factory.Close()

	// Alias call
	_, err = factory.GetLightningClient().SetAlias(context.Background(), &lnrpc.SetAliasRequest{Alias: alias})
	if err != nil {
		err := exception.NewError("cannot change the alias", err, exception.NewExampleError)
		return err
	}

	// Color call
	_, err = factory.GetLightningClient().SetColor(context.Background(), &lnrpc.SetColorRequest{Color: color})
	if err != nil {
		err := exception.NewError("cannot change the color", err, exception.NewExampleError)
		return err
	}

	return nil
}
