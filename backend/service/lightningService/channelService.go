package lightningService

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	exception "github.com/Lmare/lightning-playground/backend/exception"
	lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	routerrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
	streamService "github.com/Lmare/lightning-playground/backend/service/streamService"
)

// Connect to a new Pair
func AddPeer(dataClient nodeModel.LndClientAuthData, uri string) error {
	client, conn, err := getLightningClient(dataClient)
	if err != nil {
		return exception.NewError("cannot init Lightning Client", err, exception.NewExampleError)
	}
	defer conn.Close()

	parts := strings.Split(uri, "@")
	_, err = client.ConnectPeer(context.Background(), &lnrpc.ConnectPeerRequest{Addr: &lnrpc.LightningAddress{Pubkey: parts[0], Host: parts[1]}, Perm: true})
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
		NodePubkey:         pubkeyBytes,
		LocalFundingAmount: amount,
		SatPerVbyte:        1, // fee on chain
		UseBaseFee:         true,
		BaseFee:            1000, // 1000 = 1 sat
		UseFeeRate:         true,
		FeeRate:            1000, // nb stat by 1 000 000 of sat
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
func MakePaiment(dataClient nodeModel.LndClientAuthData, paymentRequest string) error {
	client, conn, err := getRouterClient(dataClient)
	if err != nil {
		return exception.NewError("Cannot init Router Client", err, exception.NewExampleError)
	}

	feeLimitMsat := estimateFeeLimitMsat(client, paymentRequest)

	// send the payment with an explicit fee limit so routes with non-zero fees are considered
	stream, err := client.SendPaymentV2(context.Background(), &routerrpc.SendPaymentRequest{
		PaymentRequest: paymentRequest,
		FeeLimitMsat:   feeLimitMsat,
	})
	if err != nil {
		return exception.NewError("Error on sending payment", err, exception.NewExampleError)
	}
	streamService.StreamResult(streamService.StreamWrapper[lnrpc.Payment]{
		RecvCallback:  stream.Recv,
		CloseCallback: conn.Close,
	})

	return nil
}

// estimateFeeLimitMsat calls EstimateRouteFee and applies a margin/fallback
func estimateFeeLimitMsat(client routerrpc.RouterClient, paymentRequest string) int64 {
	feeResp, err := client.EstimateRouteFee(context.Background(), &routerrpc.RouteFeeRequest{PaymentRequest: paymentRequest})
	if err != nil {
		fmt.Println("EstimateRouteFee failed, using default fee limit (100 sat):", err)
		return int64(100000) // 100 sat fallback
	}

	// add a margin (20%) to the estimated routing fee
	feeLimitMsat := feeResp.RoutingFeeMsat + (feeResp.RoutingFeeMsat / 5)
	if feeLimitMsat == 0 {
		feeLimitMsat = int64(100000)
	}
	return feeLimitMsat
}
