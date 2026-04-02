package testsupport

import (
	"net/http"

	config "github.com/Lmare/lightning-playground"
	"github.com/Lmare/lightning-playground/backend/internal/bootstrap"
	lndgrpc "github.com/Lmare/lightning-playground/backend/internal/lightning/adapter/grpc"
)

// TestHarness holds components for testing the application, such as the router and mock repositories.
type TestHarness struct {
	Router                http.Handler
	MockUserRepo          *MockUserRepository
	MockNodeRepo          *MockNodeRepository
	fakeGrpcClientFactory *lndgrpc.GrpcClientFactoryMock
}

// NewTestHarness initializes the test harness with default mock implementations.
func NewTestHarness() *TestHarness {
	// random test configuration
	cfg := &config.Config{
		ProjectPath: "testdata",
		NodeStorage: "testdata/nodes",
		Version:     "1.0.0",
	}

	// mock repositories
	mockuserRepo := &MockUserRepository{}
	mockNodeRepo := &MockNodeRepository{}

	// mock factories
	grpcClientFactoryMock := lndgrpc.NewGrpcClientFactoryMock(nil, nil) // TODO

	// wiring with mocks
	repos := &bootstrap.Repositories{
		User: mockuserRepo,
		Node: mockNodeRepo,
	}

	factories := &bootstrap.Factories{
		GrpcClientFactory: grpcClientFactoryMock,
	}

	services := bootstrap.InitServices(repos, factories, cfg)
	handlers := bootstrap.InitHandlers(services, cfg)
	router := bootstrap.InitRouter(handlers)

	return &TestHarness{
		Router:                router,
		MockUserRepo:          mockuserRepo,
		MockNodeRepo:          mockNodeRepo,
		fakeGrpcClientFactory: grpcClientFactoryMock,
	}
}
