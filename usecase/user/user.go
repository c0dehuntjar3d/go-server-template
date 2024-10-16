package user

import (
	"app/config"
	userHandler "app/internal/controller/rest/user"
	userRepo "app/internal/repository/postgres/user"
	userService "app/internal/service/user"
	"app/pkg/database"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"app/usecase"
)

const name = "user crud"

type UserUseCase struct {
	Handler    userHandler.UserHandler
	Repository userRepo.UserRepository
	Service    userService.UserService
}

var _ usecase.UseCase = (*UserUseCase)(nil)

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{}
}

func (u UserUseCase) Name() string {
	return name
}

func (u *UserUseCase) Initialize(
	log logger.Interface,
	cfg config.Config,
	server *httpserver.Server,
	db *database.Postgres,
) (usecase.UseCase, error) {
	repo, err := userRepo.NewUserRepository(db, log)
	if err != nil {
		return nil, err
	}

	service, err := userService.NewUserService(repo, log)
	if err != nil {
		return nil, err
	}

	handler, err := userHandler.NewUserHandler(service, log, server.Mux)
	if err != nil {
		return nil, err
	}

	return &UserUseCase{
		Handler:    *handler,
		Repository: repo,
		Service:    service,
	}, nil
}
