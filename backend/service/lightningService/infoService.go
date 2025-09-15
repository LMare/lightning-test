package lightningService

import (
	"context"

	lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	exception "github.com/Lmare/lightning-test/backend/exception"
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

// return node Information
func GetUsefullInfo(dataClient LndClientAuthData) (*InfoLndNode, error) {
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


func UpdateAliasAndColor(dataClient LndClientAuthData, alias string, color string) error {
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
