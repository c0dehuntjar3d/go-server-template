package initializer

import (
	"errors"
	"fmt"

	"app/internal/pkg/config"
	"app/internal/pkg/database"
	"app/internal/pkg/httpserver"
	"app/internal/pkg/logger"
)

var (
	ErrEmptyConfig = errors.New("empty configuration file")
)

type Initializer struct {
	DB     *database.Postgres
	Logger logger.Interface
	Server *httpserver.Server
}

func InitApplication() (*Initializer, *config.Config) {
	var err error
	cfg, err := config.LoadOrGetSingleton()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	initialize, err := New(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize: %w", err))
	}
	initialize.Logger.Info("Configuration was loaded success")

	return initialize, cfg
}

func DefaultApplication() *Initializer {
	cfg := config.Default()

	log, err := logger.NewZap(cfg.Log)
	if err != nil {
		fmt.Println(err)
	}

	server := httpserver.New(cfg.Http, cfg.App, log)

	return &Initializer{
		DB:     nil,
		Logger: log,
		Server: server,
	}
}

func New(cfg *config.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}

	log, err := logger.NewZap(cfg.Log)
	if err != nil {
		return nil, err
	}

	db, err := database.NewOrGetSingletonPostgres(cfg.DB, log)
	if err != nil {
		return nil, err
	}

	server := httpserver.New(cfg.Http, cfg.App, log)

	return &Initializer{
		DB:     db,
		Logger: log,
		Server: server,
	}, nil
}
