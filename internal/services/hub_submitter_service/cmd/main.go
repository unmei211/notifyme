package main

import (
	httpserver "github.com/unmei211/notifyme/internal/pkg/http_server/server"
	httpshudown "github.com/unmei211/notifyme/internal/pkg/http_server/shutdown"
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	inboxprovider "github.com/unmei211/notifyme/internal/pkg/inbox/provider"
	pginbox "github.com/unmei211/notifyme/internal/pkg/inbox/repository"
	"github.com/unmei211/notifyme/internal/pkg/kafka"
	"github.com/unmei211/notifyme/internal/pkg/logger"
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
				config.ProvideLoggerConfig,
				config.ProvideDatabaseConfig,
				config.ProvideInboxConfig,
				config.ProvideKafkaConfig,
				// Logger
				logger.InitLogger,
				// Context
				httpshudown.NewContext,
				// Database
				orm.InitGorm,
				// Database - Repository
				pginbox.InitRepository,
				// Inbox
				pginbox.PutHandlerInboxProvider,
				inbox.InitInbox,
				inboxprovider.ConsumeHandlerInboxProvider,
				// Messaging
				kafka.Init,
				// Server
				httpserver.NewHttpServer,
			),
			fx.Invoke(orm.Migrate),
			fx.Invoke(server.RunServers),
			fx.Invoke(kafka.LaunchConsumers),
		),
	).Run()
}
