package bootstrap

import (
	"database/sql"
	"net/http"

	config "github.com/Lmare/lightning-playground"
	"github.com/Lmare/lightning-playground/backend/internal/httpapi"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/adapter/filesystem"
	lndgrpc "github.com/Lmare/lightning-playground/backend/internal/lightning/adapter/grpc"
	lightninghttp "github.com/Lmare/lightning-playground/backend/internal/lightning/adapter/http"
	lightningapp "github.com/Lmare/lightning-playground/backend/internal/lightning/application"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/port"
	"github.com/Lmare/lightning-playground/backend/internal/platform/postgres"
	exception "github.com/Lmare/lightning-playground/backend/internal/shared/exception"
	userhttp "github.com/Lmare/lightning-playground/backend/internal/user/adapter/http"
	userpostgres "github.com/Lmare/lightning-playground/backend/internal/user/adapter/postgres"
	userapp "github.com/Lmare/lightning-playground/backend/internal/user/application"
	userport "github.com/Lmare/lightning-playground/backend/internal/user/port"
)

// Repositories groups outbound ports wired for the application.
type Repositories struct {
	User userport.UserRepository
	Node port.NodeRepository
}

// Factories groups technical factories (e.g. LND gRPC).
type Factories struct {
	GrpcClientFactory lndgrpc.GrpcClientFactory
}

type services struct {
	user          *userapp.UserService
	node          *lightningapp.NodeService
	lightningInfo *lndgrpc.LightningInfoService
	channel       *lndgrpc.ChannelService
}

type handlers struct {
	user        *userhttp.Handler
	lightning   *lightninghttp.Handler
	probe       *httpapi.ProbeHandler
	version     *httpapi.VersionHandler
	streamEvent *httpapi.StreamHandler
}

// ------------------- Initialization functions -------------------

// InitDB opens the Postgres connection.
func InitDB(cfg *config.Config) (*sql.DB, error) {
	return postgres.NewPostgresDB(cfg.DSN)
}

// InitRepositories builds persistence adapters.
func InitRepositories(db *sql.DB, conf *config.Config) *Repositories {
	return &Repositories{
		User: userpostgres.NewPostgresUserRepository(db),
		Node: filesystem.NewNodeRepositoryFileSystem(conf),
	}
}

// InitFactories builds shared factories.
func InitFactories() *Factories {
	return &Factories{
		GrpcClientFactory: lndgrpc.NewGrpcClientFactory(),
	}
}

// InitServices wires the application layer and LND services.
func InitServices(r *Repositories, f *Factories, conf *config.Config) *services {
	return &services{
		user:          userapp.NewUserService(r.User),
		node:          lightningapp.NewNodeService(*conf, r.Node),
		lightningInfo: lndgrpc.NewLightningInfoService(f.GrpcClientFactory),
		channel:       lndgrpc.NewChannelService(f.GrpcClientFactory),
	}
}

// InitHandlers builds HTTP adapters.
func InitHandlers(s *services, conf *config.Config) *handlers {
	return &handlers{
		user:        userhttp.NewHandler(s.user),
		lightning:   lightninghttp.NewHandler(s.lightningInfo, s.channel, s.node),
		probe:       httpapi.NewProbeHandler(),
		version:     httpapi.NewVersionHandler(conf.Version),
		streamEvent: httpapi.NewStreamHandler(),
	}
}

// InitApp is the root composition entrypoint for the backend binary.
func InitApp(cfg *config.Config) (*sql.DB, http.Handler, error) {
	db, err := InitDB(cfg)
	if err != nil {
		return nil, nil, exception.NewError("Failed to initialize database", err, exception.NewExampleError)
	}
	repos := InitRepositories(db, cfg)
	factories := InitFactories()
	services := InitServices(repos, factories, cfg)
	handlers := InitHandlers(services, cfg)
	return db, InitRouter(handlers), nil
}
