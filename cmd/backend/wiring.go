package main

import (
	"database/sql"

	config "github.com/Lmare/lightning-playground"
	handler "github.com/Lmare/lightning-playground/backend/handler"
	db "github.com/Lmare/lightning-playground/backend/infrastructure/db"
	repository "github.com/Lmare/lightning-playground/backend/repository"
	lightningService "github.com/Lmare/lightning-playground/backend/service/lightningService"
	nodeService "github.com/Lmare/lightning-playground/backend/service/nodeService"
	streamService "github.com/Lmare/lightning-playground/backend/service/streamService"
	userService "github.com/Lmare/lightning-playground/backend/service/userService"
)

type repositories struct {
	user repository.UserRepository
}

type factories struct {
	grpcClientFactory lightningService.GrpcClientFactory
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

func initDB(cfg *config.Config) (*sql.DB, error) {
	return db.NewPostgresDB(cfg.DSN)
}

func initRepositories(db *sql.DB) *repositories {
	return &repositories{
		user: repository.NewPostgresUserRepository(db),
	}
}

func initFactories() *factories {
	return &factories{
		grpcClientFactory: lightningService.NewGrpcClientFactory(),
	}
}

func initServices(r *repositories, f *factories, conf *config.Config) *services {
	return &services{
		user:          userService.NewUserService(r.user),
		node:          nodeService.NewNodeService(*conf),
		stream:        streamService.NewStreamService(),
		lightningInfo: lightningService.NewLightningInfoService(f.grpcClientFactory),
		channel:       lightningService.NewChannelService(f.grpcClientFactory),
	}
}

func initHandlers(s *services, conf *config.Config) *handlers {
	return &handlers{
		user:        handler.NewUserHandler(s.user),
		lightning:   handler.NewLightningHandler(s.lightningInfo, s.channel, s.node),
		probe:       handler.NewProbeHandler(),
		version:     handler.NewVersionHandler(conf.Version),
		streamEvent: handler.NewStreamEventHandler(s.stream),
	}
}
