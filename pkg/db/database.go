package db

import (
	"context"
	"fmt"
	"go-server/config"
	"go-server/pkg/logger"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DefaultConnAttempts = 10
	DefaultConnTimeout  = time.Second
)

type Database struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool

	logger logger.Interface
}

var pg *Database
var hdlOnce sync.Once

func NewOrGetSingleton(cfg *config.DB, logger logger.Interface) *Database {
	if cfg == nil || cfg.URL == "" {
		return nil
	}

	hdlOnce.Do(func() {
		db, err := newDatabase(cfg, logger)
		if err != nil {
			panic(err)
		}

		pg = db
	})

	return pg
}

func newDatabase(cfg *config.DB, log logger.Interface) (*Database, error) {
	pg = &Database{
		logger:       log,
		maxPoolSize:  cfg.PoolMax,
		connAttempts: cfg.ConnectionAttempts,
		connTimeout:  time.Duration(cfg.ConnectionTimeout),
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Info(fmt.Sprint("Postgres is trying to connect, attempts left: ", pg.connAttempts))

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Database) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
