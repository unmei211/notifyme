package main

import (
	"context"

	messaginginbox "github.com/unmei211/notifyme/internal/pkg/adapters/messaging/inbox"
	httpserver "github.com/unmei211/notifyme/internal/pkg/http_server/server"
	httpshudown "github.com/unmei211/notifyme/internal/pkg/http_server/shutdown"
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/kafka"
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/orm"
	"github.com/unmei211/notifyme/internal/services/hub_submitter/config"
	"github.com/unmei211/notifyme/internal/services/hub_submitter/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				// Config
				config.InitConfig,
				func(config *config.Config) (
					*httpserver.Config,
					*logger.Config,
					*orm.DatabaseConfig,
					*kafka.Config,
					*messaging.Config,
					*inbox.Config,
				) {
					return config.HttpServer,
						config.Logger,
						config.Database,
						config.Kafka,
						config.Messaging,
						config.Inbox
				},
				// Logger
				logger.InitLogger,
				// Context
				httpshudown.NewContext,
				// Database
				orm.InitGorm,
				// MessageBoxing - Repository
				inbox.InitRepository,
				// MessageBoxing - Service
				inbox.InitService,
				messaginginbox.InitHandlers,
				// MessageBoxing
				inbox.InitMessageBoxing,
				messaginginbox.InitConsumer,
				// Messaging
				fx.Annotate(
					kafka.Init,
					fx.ResultTags(`name:"kafka.manager.producer"`, `name:"kafka.manager.fetcher"`),
				),
				//kafka.Init,
				// Server
				httpserver.NewHttpServer,
			),
			fx.Invoke(orm.Migrate),
			fx.Invoke(server.RunServers),
			fx.Invoke(
				fx.Annotate(
					func(manager messaging.IFetcherManager, ctx context.Context) {
						kafka.LaunchFetcher(manager, ctx)
					},
					fx.ParamTags(`name:"kafka.manager.fetcher"`),
				),
			),
		),
	).Run()
}
