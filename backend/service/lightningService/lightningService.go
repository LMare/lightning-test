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

func Test() {
    // Chemins vers les fichiers
	basePath := "/home/louis/Documents/Dev/lightning-test/nodes-storage/lightning-test_lnd1_1"
    tlsCertPath := basePath + "/cert/tls.cert"
    macaroonPath := basePath + "/macaroons/admin.macaroon"
    lndAddress := "localhost:10009"

    // Charger le certificat TLS
    cert, err := ioutil.ReadFile(tlsCertPath)
    if err != nil {
        log.Fatalf("Erreur lecture TLS cert: %v", err)
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    creds := credentials.NewClientTLSFromCert(certPool, "")

    // Charger le macaroon
    macaroonBytes, err := ioutil.ReadFile(macaroonPath)
    if err != nil {
        log.Fatalf("Erreur lecture macaroon: %v", err)
    }
    macaroonHex := fmt.Sprintf("%x", macaroonBytes)

    // Créer un dial gRPC sécurisé
    conn, err := grpc.Dial(
        lndAddress,
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(macaroonCreds{macaroonHex}),
    )
    if err != nil {
        log.Fatalf("Erreur connexion gRPC: %v", err)
    }
    defer conn.Close()

    client := lnrpc.NewLightningClient(conn)

    // Exemple d'appel : GetInfo
    resp, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
    if err != nil {
        log.Fatalf("Erreur GetInfo: %v", err)
    }

    fmt.Printf("Node alias: %s\n", resp.Alias)
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
