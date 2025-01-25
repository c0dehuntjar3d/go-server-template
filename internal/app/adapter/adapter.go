package adapter

import (
	"app/internal/pkg/config"
	"app/internal/pkg/database"
	"app/internal/pkg/httpserver"
	"app/internal/pkg/logger"
)

type Adapter interface {
	Initialize(
		log logger.Interface,
		cfg config.Config,
		server *httpserver.Server,
		db *database.Postgres,
	) error

	Name() string
}
