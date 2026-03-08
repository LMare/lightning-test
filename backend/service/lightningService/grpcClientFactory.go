package lightningService

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"

	exception "github.com/Lmare/lightning-playground/backend/exception"
	lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	routerrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	nodeModel "github.com/Lmare/lightning-playground/backend/model/nodeModel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcClientFactory interface {
	newClient(dataClient nodeModel.LndClientAuthData) (grpcClient, error)
}

type grpcClient interface {
	GetLightningClient() lnrpc.LightningClient
	GetRouterClient() routerrpc.RouterClient
	Close() error
}

//----------------- GrpcClient Factory -----------------

type grpcClientFactoryImpl struct {
}

// NewGrpcClientFactory crée une factory avec une connexion réelle
func (f *grpcClientFactoryImpl) newClient(dataClient nodeModel.LndClientAuthData) (grpcClient, error) {
	conn, err := f.createGrpcClientConn(dataClient)
	if err != nil {
		err := exception.NewError("Erreur création canal de communication", err, exception.NewExampleError)
		return nil, err
	}

	return &grpcClientImpl{
		lightningClient: lnrpc.NewLightningClient(conn),
		routerClient:    routerrpc.NewRouterClient(conn),
		conn:            conn,
	}, nil
}

// createGrpcClientConn create a secure gRPC connection
func (f *grpcClientFactoryImpl) createGrpcClientConn(dataClient nodeModel.LndClientAuthData) (*grpc.ClientConn, error) {

	// Charger le certificat TLS
	cert, err := os.ReadFile(dataClient.TlsCertPath)
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture TLS cert : %s", dataClient.TlsCertPath), err, exception.NewExampleError)
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	creds := credentials.NewClientTLSFromCert(certPool, "")

	// Charger le macaroon
	macaroonBytes, err := os.ReadFile(dataClient.MacaroonPath)
	if err != nil {
		err := exception.NewError(fmt.Sprintf("Erreur lecture macaroon : %s", dataClient.MacaroonPath), err, exception.NewExampleError)
		return nil, err
	}
	macaroonHex := fmt.Sprintf("%x", macaroonBytes)

	// Créer un dial gRPC sécurisé
	return grpc.NewClient(
		dataClient.LndAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMacaroonCreds(macaroonHex)),
	)
}

// constructor for the factory
func NewGrpcClientFactory() GrpcClientFactory {
	return &grpcClientFactoryImpl{}
}

//----------------- Macaroon Credentials -----------------

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

func newMacaroonCreds(macaroonHex string) credentials.PerRPCCredentials {
	return macaroonCreds{macaroon: macaroonHex}
}

//----------------- GrpcClient Implementation -----------------

// GrpcClientFactory encapsule les clients gRPC pour faciliter les tests
type grpcClientImpl struct {
	lightningClient lnrpc.LightningClient
	routerClient    routerrpc.RouterClient
	conn            *grpc.ClientConn
}

// GetLightningClient retourn le Lightning client
func (f *grpcClientImpl) GetLightningClient() lnrpc.LightningClient {
	return f.lightningClient
}

// GetRouterClient retourn the Router client
func (f *grpcClientImpl) GetRouterClient() routerrpc.RouterClient {
	return f.routerClient
}

// Close the gRPC connection
func (f *grpcClientImpl) Close() error {
	if f.conn != nil {
		return f.conn.Close()
	}
	return nil
}

// ---------------- Mock for test ----------------

// Mock Factory
type GrpcClientFactoryMock struct {
	lightningClient lnrpc.LightningClient
	routerClient    routerrpc.RouterClient
}

// create a mocked grpc client
func (f *GrpcClientFactoryMock) newClient(dataClient nodeModel.LndClientAuthData) (grpcClient, error) {
	return &grpcClientImpl{
		lightningClient: f.lightningClient,
		routerClient:    f.routerClient,
		conn:            nil,
	}, nil
}

// constructor for the mock factory
func NewGrpcClientFactoryMock(lightningClient lnrpc.LightningClient, routerClient routerrpc.RouterClient) *GrpcClientFactoryMock {
	return &GrpcClientFactoryMock{
		lightningClient: lightningClient,
		routerClient:    routerClient,
	}
}

// assertion compile-time
var _ GrpcClientFactory = (*GrpcClientFactoryMock)(nil)
