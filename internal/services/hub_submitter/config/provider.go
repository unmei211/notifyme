package config

import (
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"github.com/unmei211/notifyme/internal/pkg/orm"
)

func ProvideLoggerConfig(
	config *Config,
) *logger.Config {
	return config.Logger
}

func ProvideDatabaseConfig(
	config *Config) *orm.DatabaseConfig {
	logger.Log.Debugf("Provide database config")
	return config.Database
}
