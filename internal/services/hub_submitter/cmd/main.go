package main

import (
	httpserver "github.com/unmei211/notifyme/internal/pkg/http_server/server"
	httpshudown "github.com/unmei211/notifyme/internal/pkg/http_server/shutdown"
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
				config.InitConfig,
				config.ProvideLoggerConfig,
				logger.InitLogger,
				httpshudown.NewContext,
				config.ProvideDatabaseConfig,
				orm.InitGorm,
				kafka.Init,
				httpserver.NewHttpServer,
			),
			fx.Invoke(orm.Migrate),
			fx.Invoke(server.RunServers),
		),
	).Run()
}
