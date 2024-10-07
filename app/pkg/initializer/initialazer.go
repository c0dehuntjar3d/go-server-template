package initializer

import (
	"app/config"
	"app/pkg/db"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"errors"
	"fmt"
)

var (
	ErrEmptyConfig = errors.New("empty configuration file")
)

type Initializer struct {
	DB     *db.Database
	Logger logger.Interface
	Server *httpserver.Server
}

func InitApplicaiton() (*Initializer, *config.Config) {
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

	db, err := db.NewOrGetSingleton(cfg.DB, log)
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
