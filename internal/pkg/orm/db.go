package orm

import (
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/unmei211/notifyme/internal/pkg/logger"
	"go.uber.org/zap"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Url      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	DB       string `mapstructure:"db"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
}

type Gorm struct {
	DB     *gorm.DB
	log    *zap.SugaredLogger
	config *DatabaseConfig
}

func buildConnectionString(cfg *DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&search_path=%s",
		cfg.User,
		cfg.Password,
		cfg.Url,
		cfg.DB,
		cfg.Schema,
	)
}

func InitGorm(cfg *DatabaseConfig, log *zap.SugaredLogger) (*Gorm, error) {
	connectionStr := buildConnectionString(cfg)
	db, err := gorm.Open(gormpostgres.Open(connectionStr), &gorm.Config{})

	if err != nil {
		log.Errorf("Failed to initialize db")
		return nil, err
	}

	return &Gorm{
		config: cfg,
		DB:     db,
		log:    log,
	}, nil
}

func Migrate(cfg *DatabaseConfig) error {
	m, err := migrate.New(
		"file://./migrations",
		buildConnectionString(cfg),
	)
	if err != nil {
		logger.Log.Errorf("Failed to initialize migrator")
		return err
	}
	err = m.Up()
	if err != nil {
		logger.Log.Errorf("Failed to process migrations")
		return err
	}
	return nil
}
