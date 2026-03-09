package testsupport

import (
	"net/http"

	config "github.com/Lmare/lightning-playground"
	app "github.com/Lmare/lightning-playground/backend/app"
	lightningService "github.com/Lmare/lightning-playground/backend/service/lightningService"
)

// TestHarness is a struct that holds all the necessary components for testing the application, such as the router and mock repositories.
type TestHarness struct {
	Router                http.Handler
	MockUserRepo          *MockUserRepository
	MockNodeRepo          *MockNodeRepository
	fakeGrpcClientFactory *lightningService.GrpcClientFactoryMock
}

// NewTestHarness initializes the test harness with default mock implementations for repositories and factories.
func NewTestHarness() *TestHarness {
	// random test configuration
	config := &config.Config{
		ProjectPath: "testdata",
		NodeStorage: "testdata/nodes",
		Version:     "1.0.0",
	}

	// mock repositories

	mockuserRepo := &MockUserRepository{}

	mockNodeRepo := &MockNodeRepository{}

	// mock factories

	grpcClientFactoryMock := lightningService.NewGrpcClientFactoryMock(nil, nil) // TODO

	// wiring with mocks

	repos := &app.Repositories{
		User: mockuserRepo,
		Node: mockNodeRepo,
	}

	factories := &app.Factories{
		GrpcClientFactory: grpcClientFactoryMock,
	}

	services := app.InitServices(repos, factories, config)
	handlers := app.InitHandlers(services, config)
	router := app.InitRouter(handlers)

	return &TestHarness{
		Router:                router,
		MockUserRepo:          mockuserRepo,
		MockNodeRepo:          mockNodeRepo,
		fakeGrpcClientFactory: grpcClientFactoryMock,
	}

}
