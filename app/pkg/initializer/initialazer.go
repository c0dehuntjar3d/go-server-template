package initializer

import (
	"app/config"
	"app/pkg/database"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"app/usecase"
	"app/usecase/user"
	"errors"
	"fmt"
)

var (
	ErrEmptyConfig = errors.New("empty configuration file")
)

type Initializer struct {
	DB     *database.Postgres
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

	loadUseCases(initialize, *cfg)
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

func loadUseCases(initializr *Initializer, cfg config.Config) {
	usecases := []usecase.UseCase{
		user.NewUserUseCase(),
	}

	for _, uc := range usecases {
		_, err := uc.Initialize(initializr.Logger, cfg, initializr.Server, initializr.DB)
		if err != nil {
			initializr.Logger.Fatal(fmt.Sprintf("Usecase %s load error: %s", uc.Name(), err.Error()))
		} else {
			initializr.Logger.Debug(fmt.Sprintf("Usecase %s load success", uc.Name()))
		}

	}
}
