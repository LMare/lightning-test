package lightningService


import (
	"context"
	"strings"
	"encoding/hex"

	routerrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	exception "github.com/Lmare/lightning-test/backend/exception"
	nodeModel "github.com/Lmare/lightning-test/backend/model/nodeModel"
	streamService "github.com/Lmare/lightning-test/backend/service/streamService"
)


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



// Open Channel
func OpenChannel(dataClient nodeModel.LndClientAuthData, pubkeyHex string, amount int64) error {

	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		return exception.NewError("Cannot init Lightning Client", err, exception.NewExampleError)
    }
    defer conn.Close()

	pubkeyBytes, err := hex.DecodeString(pubkeyHex)
	if err != nil {
	    return exception.NewError("Error on decode pubkeyHex", err, exception.NewExampleError)
	}

	r := lnrpc.OpenChannelRequest{
		NodePubkey: pubkeyBytes,
		LocalFundingAmount: amount,
		SatPerVbyte: 1, // fee on chain
		UseBaseFee: true,
		BaseFee: 1000, // 1000 = 1 sat
		UseFeeRate: true,
		FeeRate: 1000, // nb stat by 1 000 000 of sat
	}

	// TODO @GoodToHave utiliser le retour de OpenChannel pour faire des pipes de notifications sur l'IHM
	_, err = client.OpenChannelSync(context.Background(), &r)
	if err != nil {
		return exception.NewError("Error on opening a channel", err, exception.NewExampleError)
	}
	return nil
}



// create an invoice which while expire in 5min
// return the payment request
func CreateQuickInvoice(dataClient nodeModel.LndClientAuthData, memo string, amount int64) (string, error) {

	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		return "", exception.NewError("Cannot init Lightning Client", err, exception.NewExampleError)
    }
    defer conn.Close()

	fiveMin := int64(300)

	i, err := client.AddInvoice(context.Background(), &lnrpc.Invoice{Expiry: fiveMin, Memo: memo, Value: amount})
	if err != nil {
		return "", exception.NewError("Error on creating invoice", err, exception.NewExampleError)
	}

	return i.PaymentRequest, nil
}

// pay the invoice
// return streamId, error
func MakePaiment(dataClient nodeModel.LndClientAuthData, paymentRequest string) (string, error) {
	client, conn, err := getRouterClient(dataClient)
	if err != nil {
		return "", exception.NewError("Cannot init Router Client", err, exception.NewExampleError)
    }

	stream, err := client.SendPaymentV2(context.Background(), &routerrpc.SendPaymentRequest{PaymentRequest: paymentRequest})
	if err != nil {
		return "", exception.NewError("Error on creating invoice", err, exception.NewExampleError)
	}
	streamId := streamService.KeepStream(streamService.StreamWrapper[lnrpc.Payment]{
		RecvCallback: stream.Recv,
		CloseCallback: conn.Close,
	})

	return streamId, nil
}
