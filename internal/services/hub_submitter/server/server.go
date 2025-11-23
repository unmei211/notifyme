package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	httpserver "github.com/unmei211/notifyme/internal/pkg/http_server/server"
	"github.com/unmei211/notifyme/internal/services/hub_submitter/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"net/http"
)

func RunServers(
	lc fx.Lifecycle,
	log *zap.SugaredLogger,
	e *echo.Echo,
	//grpcServer *grpc.GrpcServer,
	ctx context.Context,
	cfg *config.Config,
) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := httpserver.RunHttpServer(ctx, e, log, cfg.HttpServer); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			//go func() {
			//	if err := grpcServer.RunGrpcServer(ctx); !errors.Is(err, http.ErrServerClosed) {
			//		log.Fatalf("error running grpc server: %v", err)
			//	}
			//}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, config.GetMicroserviceName(cfg.ServiceName))
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
