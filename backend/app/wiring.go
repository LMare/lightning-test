package app

import (
	"database/sql"
	"net/http"

	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/exception"
	handler "github.com/Lmare/lightning-playground/backend/handler"
	db "github.com/Lmare/lightning-playground/backend/infrastructure/db"
	repository "github.com/Lmare/lightning-playground/backend/repository"
	lightningService "github.com/Lmare/lightning-playground/backend/service/lightningService"
	nodeService "github.com/Lmare/lightning-playground/backend/service/nodeService"
	streamService "github.com/Lmare/lightning-playground/backend/service/streamService"
	userService "github.com/Lmare/lightning-playground/backend/service/userService"
)

type Repositories struct {
	User repository.UserRepository
	Node repository.NodeRepository
}

type Factories struct {
	GrpcClientFactory lightningService.GrpcClientFactory
}

type services struct {
	user          *userService.UserService
	node          *nodeService.NodeService
	stream        *streamService.StreamService
	lightningInfo *lightningService.LightningInfoService
	channel       *lightningService.ChannelService
}

type handlers struct {
	user        *handler.UserHandler
	lightning   *handler.LightningHandler
	probe       *handler.ProbeHandler
	version     *handler.VersionHandler
	streamEvent *handler.StreamEventHandler
}

// ------------------- Initialization functions -------------------

func InitDB(cfg *config.Config) (*sql.DB, error) {
	return db.NewPostgresDB(cfg.DSN)
}

func InitRepositories(db *sql.DB, conf *config.Config) *Repositories {
	return &Repositories{
		User: repository.NewPostgresUserRepository(db),
		Node: repository.NewNodeRepositoryFileSystem(conf),
	}
}

func InitFactories() *Factories {
	return &Factories{
		GrpcClientFactory: lightningService.NewGrpcClientFactory(),
	}
}

func InitServices(r *Repositories, f *Factories, conf *config.Config) *services {
	return &services{
		user:          userService.NewUserService(r.User),
		node:          nodeService.NewNodeService(*conf, r.Node),
		stream:        streamService.NewStreamService(),
		lightningInfo: lightningService.NewLightningInfoService(f.GrpcClientFactory),
		channel:       lightningService.NewChannelService(f.GrpcClientFactory),
	}
}

func InitHandlers(s *services, conf *config.Config) *handlers {
	return &handlers{
		user:        handler.NewUserHandler(s.user),
		lightning:   handler.NewLightningHandler(s.lightningInfo, s.channel, s.node),
		probe:       handler.NewProbeHandler(),
		version:     handler.NewVersionHandler(conf.Version),
		streamEvent: handler.NewStreamEventHandler(s.stream),
	}
}

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
