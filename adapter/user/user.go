package user

import (
	"app/adapter"
	userHandler "app/internal/controller/rest/user"
	userRepo "app/internal/repository/postgres/user"
	userCase "app/internal/usecase/user"
	"app/pkg/config"
	"app/pkg/database"
	"app/pkg/httpserver"
	"app/pkg/logger"
)

const name = "user crud"

type UserAdapter struct {
	Handler    userHandler.UserHandler
	Repository userRepo.UserRepository
	Service    userCase.UserService
}

var _ adapter.Adapter = (*UserAdapter)(nil)

func NewUserAdapter() *UserAdapter {
	return &UserAdapter{}
}

func (u UserAdapter) Name() string {
	return name
}

func (u *UserAdapter) Initialize(
	log logger.Interface,
	cfg config.Config,
	server *httpserver.Server,
	db *database.Postgres,
) error {
	repository, err := userRepo.NewUserRepository(db, log)
	if err != nil {
		return err
	}

	service, err := userCase.NewUserService(repository, log)
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
