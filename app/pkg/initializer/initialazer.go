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

func WithDefault() *Initializer {
	cfg := config.WithDefault()

	log, err := logger.NewZap(cfg.Log)
	if err != nil {
		fmt.Println(err)
	}

	db := db.NewOrGetSingleton(cfg.DB, log)
	server := httpserver.New(cfg.Http, log)

	return &Initializer{
		DB:     db,
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
	db := db.NewOrGetSingleton(cfg.DB, log)
	server := httpserver.New(cfg.Http, log)

	return &Initializer{
		DB:     db,
		Logger: log,
		Server: server,
	}, nil
}
