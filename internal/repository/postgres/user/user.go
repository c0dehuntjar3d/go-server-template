package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/domain"
	"app/internal/repository/postgres/user/converter"
	repoUser "app/internal/repository/postgres/user/model"
	"app/pkg/database"
	"app/pkg/logger"
	"app/pkg/types"
)

func NewRepository() *repository {
	return &repository{}
}

var _ UserRepository = (*repository)(nil)

type repository struct {
	db     *database.Postgres
	logger logger.Interface
}

func NewUserRepository(db *database.Postgres, logger logger.Interface) (*repository, error) {
	if db == nil {
		return nil, errors.New("db is null")
	}

	if logger == nil {
		return nil, errors.New("logger is null")
	}

	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *repository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u == nil {
		return nil, errors.New("user is null")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	r.logger.Debug(fmt.Sprintf("txID: %s [repository], creating user: %v", txID, u))

	query := `INSERT INTO users (uuid, login, password, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING uuid`

	u.CreatedAt = time.Now()
	repoUser := converter.ToUserFromDomain(u)

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		repoUser.Uuid,
		repoUser.Login,
		repoUser.Password,
		repoUser.CreatedAt,
	).Scan(&repoUser.Uuid)

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s [repository], error creating user: %v", txID, err))
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("txID: %s [repository], successfully created user with UUID: %s", txID, repoUser.Uuid))

	return converter.ToUserFromRepository(repoUser), nil
}

func (r *repository) Delete(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s [repository], deleting user with UUID: %s", txID, uuid))

	query := `UPDATE users 
              SET deleted_at = $1 
              WHERE uuid = $2 
              RETURNING uuid, login, password, created_at, updated_at, deleted_at`

	now := time.Now()
	user := &repoUser.User{}

	err := r.db.Pool.QueryRow(ctx, query, now, uuid).Scan(
		&user.Uuid,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s [repository], error soft deleting user: %v", txID, err))
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("txID: %s [repository], successfully soft deleted user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}

func (r *repository) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s [repository], fetching user with UUID: %s", txID, uuid))

	query := `SELECT uuid, login, password, created_at, updated_at, deleted_at 
              FROM users WHERE uuid = $1`

	user := &repoUser.User{}

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		uuid,
	).Scan(
		&user.Uuid,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s [repository], error fetching user: %v", txID, err))
		return nil, err
	}

	r.logger.Debug(fmt.Sprintf("txID: %s [repository], successfully fetched user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}
