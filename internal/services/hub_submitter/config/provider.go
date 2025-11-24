package config

import (
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"github.com/unmei211/notifyme/internal/pkg/orm"
	"go.uber.org/zap"
)

func ProvideLoggerConfig(
	config *Config,
) *logger.Config {
	return config.Logger
}

func ProvideDatabaseConfig(
	config *Config, log *zap.SugaredLogger) *orm.DatabaseConfig {
	log.Debugf("Provide database config")
	return config.Database
}
