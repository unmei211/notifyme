package config

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"github.com/unmei211/notifyme/internal/pkg/messaging/kafka"
	"github.com/unmei211/notifyme/internal/pkg/orm"
)

// ProvideLoggerConfig Provide logger config
func ProvideLoggerConfig(
	config *Config,
) *logger.Config {
	return config.Logger
}

// ProvideDatabaseConfig Provide database config
func ProvideDatabaseConfig(
	config *Config) *orm.DatabaseConfig {
	return config.Database
}

// ProvideKafkaConfig Provide kafka config
func ProvideKafkaConfig(
	config *Config) *kafka.Config {
	return config.Kafka
}

// ProvideInboxConfig Provide inbox config
func ProvideInboxConfig(
	config *Config,
) *inbox.Config {
	return config.Inbox
}
