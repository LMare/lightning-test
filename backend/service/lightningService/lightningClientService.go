package lightningService

import (
    "crypto/x509"
    "fmt"
    "io/ioutil"
	"context"


    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	exception "github.com/Lmare/lightning-test/backend/exception"
	nodeModel "github.com/Lmare/lightning-test/backend/model/nodeModel"
)



// get the gRPC client to interract with a node
func getLightningClient(dataClient nodeModel.LndClientAuthData) (lnrpc.LightningClient, *grpc.ClientConn, error) {

    // Charger le certificat TLS
    cert, err := ioutil.ReadFile(dataClient.TlsCertPath)
    if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture TLS cert : %s", dataClient.TlsCertPath), err, exception.NewExampleError)
		return nil, nil, err
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    creds := credentials.NewClientTLSFromCert(certPool, "")

    // Charger le macaroon
    macaroonBytes, err := ioutil.ReadFile(dataClient.MacaroonPath)
    if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture macaroon : %s", dataClient.MacaroonPath), err, exception.NewExampleError)
		return nil, nil, err
    }
    macaroonHex := fmt.Sprintf("%x", macaroonBytes)

    // Créer un dial gRPC sécurisé
    conn, err := grpc.Dial(
        dataClient.LndAddress,
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(macaroonCreds{macaroonHex}),
    )
    if err != nil {
		err := exception.NewError("Erreur création canal de communication", err, exception.NewExampleError)
		return nil, nil, err
    }

    return lnrpc.NewLightningClient(conn), conn, nil
}

// macaroonCreds permet d'ajouter le macaroon dans les métadonnées gRPC
type macaroonCreds struct {
    macaroon string
}

func (m macaroonCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
    return map[string]string{"macaroon": m.macaroon}, nil
}

func (m macaroonCreds) RequireTransportSecurity() bool {
    return true
}
