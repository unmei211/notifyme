package config

import "github.com/unmei211/notifyme/internal/pkg/logger"

func ProvideLoggerConfig(
	config *Config,
) *logger.Config {
	return config.Logger
}
