package lightningService

import (
	"context"
	"fmt"
	"strings"

	lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	exception "github.com/Lmare/lightning-test/backend/exception"
	nodeModel "github.com/Lmare/lightning-test/backend/model/nodeModel"
)

type InfoLndNode struct {
	Alias					string		`json:"alias"`
	Color					string		`json:"color"`
	NumPendingChannels		uint32		`json:"numPendingChannels"`
	NumActiveChannels		uint32		`json:"numActiveChannels"`
	NumInactiveChannels		uint32		`json:"numInactiveChannels"`
	NumPeers				uint32		`json:"numPeers"`
	BlockHeight				uint32		`json:"blockHeight"`
	Network					string		`json:"network"`
	Uris					[]string	`json:"uris"`
	SyncedToChain			bool		`json:"syncedToChain"`
	SyncedToGraph			bool		`json:"syncedToGraph"`
}

type NodeBasicInfo struct {
	Id 		int		`json:"id"`
	Alias	string	`json:"alias"`
	Color	string	`json:"color"`
}

// Get the active list of node
func GetListOfNode(descriptors []nodeModel.NodeConfigDescriptor) ([]NodeBasicInfo) {
	l := []NodeBasicInfo{}

	for _, descriptor := range descriptors {
		basicInfo, err := getBasicInfo(descriptor)
	    if err != nil {
			fmt.Println("[WARN] ", err)
			continue
	    }
		l = append(l, basicInfo)
	}
	return l
}

// get the basical info of a node
func getBasicInfo(descriptor nodeModel.NodeConfigDescriptor) (NodeBasicInfo, error) {
	client, conn, err := getLightningClient(descriptor.AuthData)
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Unable to open dial with Node[%d]", descriptor.Id), err, exception.NewExampleError)
		return NodeBasicInfo{}, err
	}
	defer conn.Close()

	resp, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Unable to getInfo of Node[%d]", descriptor.Id), err, exception.NewExampleError)
		return NodeBasicInfo{}, err
	}

	return NodeBasicInfo{Id: descriptor.Id, Alias: resp.GetAlias(), Color: resp.GetColor(),},  nil
}

// get the uri of the lnd
func GetFirstUri(dataClient nodeModel.LndClientAuthData) (string, error) {
	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		err := exception.NewError("Unable to open dial", err, exception.NewExampleError)
		return "", err
	}
	defer conn.Close()

	resp, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
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

	return uri,  nil
}


// return node Information
func GetUsefullInfo(dataClient nodeModel.LndClientAuthData) (*InfoLndNode, error) {
	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		err := exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
		return nil, err
    }
    defer conn.Close()

    resp, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
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

// Connect to a new Pair
func AddPeer(dataClient nodeModel.LndClientAuthData, uri string) error {
	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		return exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
    }
    defer conn.Close()

	parts := strings.Split(uri, "@")
	_, err = client.ConnectPeer(context.Background(), &lnrpc.ConnectPeerRequest{Addr: &lnrpc.LightningAddress{Pubkey: parts[0], Host: parts[1],},})
	if err != nil {
		return exception.NewError("Error on adding a peer", err, exception.NewExampleError)
	}
	return nil
}



// TODO:
func UpdateAliasAndColor(dataClient nodeModel.LndClientAuthData, alias string, color string) error {
	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		err := exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
		return err
    }
    defer conn.Close()

	// Todo : No Interface gRPC to do that :/
	_ = client

	return exception.NewError("Not Implemended", err, exception.NewExampleError)

}
