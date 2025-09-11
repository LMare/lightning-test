package lightningService

import (
    //"crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
	"context"


    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    //"google.golang.org/grpc/metadata"
    lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"

)

type LndClientAuthData struct {
	tlsCertPath 	string
	macaroonPath 	string
	lndAddress 		string
}

func NewLndClientAuthData(c, m, a string) LndClientAuthData {
	return LndClientAuthData{c, m, a}
}


func getLightningClient(dataClient LndClientAuthData) (lnrpc.LightningClient, *grpc.ClientConn, error) {
    // Chemins vers les fichiers


    // Charger le certificat TLS
    cert, err := ioutil.ReadFile(dataClient.tlsCertPath)
    if err != nil {
        log.Fatalf("Erreur lecture TLS cert: %v", err)
		return nil, nil, err
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    creds := credentials.NewClientTLSFromCert(certPool, "")

    // Charger le macaroon
    macaroonBytes, err := ioutil.ReadFile(dataClient.macaroonPath)
    if err != nil {
        log.Fatalf("Erreur lecture macaroon: %v", err)
		return nil, nil, err
    }
    macaroonHex := fmt.Sprintf("%x", macaroonBytes)

    // Créer un dial gRPC sécurisé
    conn, err := grpc.Dial(
        dataClient.lndAddress,
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(macaroonCreds{macaroonHex}),
    )
    if err != nil {
        log.Fatalf("Erreur connexion gRPC: %v", err)
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
