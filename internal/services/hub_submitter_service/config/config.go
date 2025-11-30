package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	httpserver "github.com/unmei211/notifyme/internal/pkg/http_server/server"
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/kafka"
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/orm"
)

type Config struct {
	ServiceName string              `mapstructure:"serviceName"`
	HttpServer  *httpserver.Config  `mapstructure:"httpServer"`
	Logger      *logger.Config      `mapstructure:"logger"`
	Database    *orm.DatabaseConfig `mapstructure:"database"`
	Kafka       *kafka.Config       `mapstructure:"kafka"`
	Messaging   *messaging.Config   `mapstructure:"messaging"`
	Inbox       *inbox.Config       `mapstructure:"inbox"`
}

func InitConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	configPath := os.Getenv("CONFIG_PATH")

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}

func GetMicroserviceName(serviceName string) string {
	return fmt.Sprintf("%s", strings.ToUpper(serviceName))
}
