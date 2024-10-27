package database

import (
	"app/pkg/config"
	"app/pkg/logger"
	"context"
	"errors"
	"fmt"

	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool

	logger logger.Interface
}

var pg *Postgres
var hdlOnce sync.Once

func NewOrGetSingletonPostgres(cfg *config.DB, logger logger.Interface) (*Postgres, error) {
	if cfg == nil || validateEmpty(cfg) {
		return nil, errors.New("cfg is not valid")
	}

	var er error
	hdlOnce.Do(func() {
		db, err := newDatabase(cfg, logger)
		if err != nil {
			er = err
		}

		pg = db
	})

	return pg, er
}

func validateEmpty(cfg *config.DB) bool {
	return cfg.Database == "" || cfg.Password == "" || cfg.User == ""
}

func newDatabase(cfg *config.DB, log logger.Interface) (*Postgres, error) {
	pg = &Postgres{
		logger:       log,
		connAttempts: cfg.ConnectionAttempts,
		connTimeout:  time.Duration(cfg.ConnectionTimeout),
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSL,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Warn(fmt.Sprint("Postgres is trying to connect, attempts left: ", pg.connAttempts))

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
