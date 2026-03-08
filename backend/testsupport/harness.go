package testsupport

import (
	"context"

	app "github.com/Lmare/lightning-playground/backend/app"
	user "github.com/Lmare/lightning-playground/backend/model/user"
	lightningService "github.com/Lmare/lightning-playground/backend/service/lightningService"
)

// TestHarness is a struct that holds all the necessary components for testing the application, such as the router and mock repositories.
type TestHarness struct {
	Router                *app.Router
	fakeUserRepo          *FakeUserRepository
	fakeNodeRepo          *FakeNodeRepository
	fakeGrpcClientFactory *lightningService.GrpcClientFactoryMock
}

// NewTestHarness initializes the test harness with default mock implementations for repositories and factories.
func NewTestHarness() *TestHarness {

	// mock repositories

	userRepo := &FakeUserRepository{
		MockFindAll: func(ctx context.Context) ([]user.UserModel, error) {
			// by default, return some dummy data for testing
			return []user.UserModel{
				{ID: "1", Nom: "Doe", Prenom: "John", Age: 30, Email: "john.doe@example.com"},
				{ID: "2", Nom: "Smith", Prenom: "Jane", Age: 25, Email: "jane.smith@example.com"},
			}, nil
		},
	}

	fakeNodeRepo := &FakeNodeRepository{
		MockGetNodesIds: func() ([]string, error) {
			// by default, return some dummy data for testing
			return []string{"node1", "node2", "node3"}, nil
		},
	}

	// mock factories

	grpcClientFactoryMock := lightningService.NewGrpcClientFactoryMock(nil, nil) // TODO

	// wiring with mocks

	repos := &app.Repositories{
		User: userRepo,
		Node: fakeNodeRepo,
	}

	factories := &app.Factories{
		GrpcClientFactory: grpcClientFactoryMock,
	}

	services := app.InitServices(repos, factories, nil)
	handlers := app.InitHandlers(services, nil)
	router := app.InitRouter(handlers)

	return &TestHarness{
		Router:                router,
		fakeUserRepo:          userRepo,
		fakeNodeRepo:          fakeNodeRepo,
		fakeGrpcClientFactory: grpcClientFactoryMock,
	}

}
