package adapter

import (
	"app/pkg/config"
	"app/pkg/database"
	"app/pkg/httpserver"
	"app/pkg/logger"
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
