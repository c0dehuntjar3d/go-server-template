package usecase

import (
	"app/config"
	"app/pkg/database"
	"app/pkg/httpserver"
	"app/pkg/logger"
)

type UseCase interface {
	Initialize(
		log logger.Interface,
		cfg config.Config,
		server *httpserver.Server,
		db *database.Postgres,
	) (UseCase, error)

	Name() string
}
