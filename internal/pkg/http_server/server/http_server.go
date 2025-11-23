package httpserver

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type Config struct {
	Port     string `mapstructure:"port" validate:"required"`
	Env      string `mapstructure:"env"`
	BasePath string `mapstructure:"basePath" validate:"required"`
	Timeout  int    `mapstructure:"timeout"`
	Host     string `mapstructure:"host"`
}

func NewHttpServer() *echo.Echo {
	e := echo.New()

	return e
}

func RunHttpServer(
	ctx context.Context,
	echo *echo.Echo,
	log *zap.SugaredLogger,
	cfg *Config) error {
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	log.Info("Try run http server")
	go func() {
		<-ctx.Done()
		err := echo.Shutdown(ctx)
		if err != nil {
			return
		}
		return
	}()

	err := echo.Start("localhost:" + cfg.Port)
	return err
}

func RegisterGroupFunc(groupName string, echo *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(echo.Group(groupName))
	return echo
}
