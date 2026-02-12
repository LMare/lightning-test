package lightningService

import (
    "crypto/x509"
    "fmt"
    "io/ioutil"
	"context"


    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	routerrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	exception "github.com/Lmare/lightning-playground/backend/exception"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
)

// GrpcClientFactory encapsule les clients gRPC pour faciliter les tests
type GrpcClientFactory struct {
	lightningClient lnrpc.LightningClient
	routerClient    routerrpc.RouterClient
	conn            *grpc.ClientConn
}

// NewGrpcClientFactory crée une factory avec une connexion réelle
func NewGrpcClientFactory(dataClient nodeModel.LndClientAuthData) (*GrpcClientFactory, error) {
	conn, err := createGrpcClientConn(dataClient)
	if err != nil {
		err := exception.NewError("Erreur création canal de communication", err, exception.NewExampleError)
		return nil, err
	}

	return &GrpcClientFactory{
		lightningClient: lnrpc.NewLightningClient(conn),
		routerClient:    routerrpc.NewRouterClient(conn),
		conn:            conn,
	}, nil
}

// NewGrpcClientFactoryWithClients crée une factory avec des clients injectés (pour les tests)
func NewGrpcClientFactoryWithClients(lightning lnrpc.LightningClient, router routerrpc.RouterClient) *GrpcClientFactory {
	return &GrpcClientFactory{
		lightningClient: lightning,
		routerClient:    router,
		conn:            nil,
	}
}

// GetLightningClient retourne le client Lightning
func (f *GrpcClientFactory) GetLightningClient() lnrpc.LightningClient {
	return f.lightningClient
}

// GetRouterClient retourne le client Router
func (f *GrpcClientFactory) GetRouterClient() routerrpc.RouterClient {
	return f.routerClient
}

// Close ferme la connexion gRPC
func (f *GrpcClientFactory) Close() error {
	if f.conn != nil {
		return f.conn.Close()
	}
	return nil
}

// createGrpcClientConn crée une connexion gRPC sécurisée
func createGrpcClientConn(dataClient nodeModel.LndClientAuthData) (*grpc.ClientConn, error) {

    // Charger le certificat TLS
    cert, err := ioutil.ReadFile(dataClient.TlsCertPath)
    if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture TLS cert : %s", dataClient.TlsCertPath), err, exception.NewExampleError)
		return nil, err
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    creds := credentials.NewClientTLSFromCert(certPool, "")

    // Charger le macaroon
    macaroonBytes, err := ioutil.ReadFile(dataClient.MacaroonPath)
    if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture macaroon : %s", dataClient.MacaroonPath), err, exception.NewExampleError)
		return nil, err
    }
    macaroonHex := fmt.Sprintf("%x", macaroonBytes)

    // Créer un dial gRPC sécurisé
    return grpc.Dial(
        dataClient.LndAddress,
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(macaroonCreds{macaroonHex}),
    )
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
