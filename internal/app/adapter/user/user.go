package user

import (
	"app/internal/app/adapter"
	userHandler "app/internal/app/controller/rest/user"
	userRepository "app/internal/app/repository/postgres/user"
	userService "app/internal/app/usecase/user"
	"app/internal/pkg/config"
	"app/internal/pkg/database"
	"app/internal/pkg/httpserver"
	"app/internal/pkg/logger"
)

const name = "user crud"

type userAdapter struct {
	Handler    userHandler.Handler
	Repository userRepository.Repository
	Service    userService.Service
}

var _ adapter.Adapter = (*userAdapter)(nil)

func NewUserAdapter() adapter.Adapter {
	return &userAdapter{}
}

func (u *userAdapter) Name() string {
	return name
}

func (u *userAdapter) Initialize(
	log logger.Interface,
	_ config.Config,
	server *httpserver.Server,
	db *database.Postgres,
) error {
	repository, err := userRepository.NewUserRepository(db, log)
	if err != nil {
		return err
	}

	service, err := userService.NewUserService(repository, log)
	if err != nil {
		return err
	}

	handler, err := userHandler.NewUserHandler(service, log, server.Mux)
	if err != nil {
		return err
	}

	u.Handler = *handler
	u.Service = service
	u.Repository = repository

	return nil
}
