package orm

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
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

func Migrate(cfg *DatabaseConfig, log *zap.SugaredLogger) error {
	migrationPath, err := filepath.Abs("/opt/service/resources/migrations")

	if err != nil {
		log.Errorf("Failed get absolute path for migrations")
		return err
	}

	log.Infof("Start init migrations. Migratios path: {%s}", migrationPath)
	m, err := migrate.New(
		"file://"+migrationPath,
		buildConnectionString(cfg),
	)
	if err != nil {
		log.Errorf("Failed to initialize migrator, err: {%+v}", err)
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No migrations to apply, database is up to date")
		} else {
			log.Errorf("Migration failed, err: %+v", err)
			return err
		}
	} else {
		log.Info("Migrations applied successfully")
	}

	return nil
}
