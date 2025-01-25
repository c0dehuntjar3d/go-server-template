package user

import (
	"app/internal/app/repository/postgres/user/converter"
	repoUser "app/internal/app/repository/postgres/user/model"
	"app/internal/domain"
	"app/internal/pkg/database"
	"app/internal/pkg/logger"
	"app/internal/pkg/types"
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, uuid string) (*domain.User, error)
	Delete(ctx context.Context, uuid string) (*domain.User, error)
}

var _ Repository = (*repository)(nil)

type repository struct {
	db     *database.Postgres
	logger logger.Interface
}

func NewUserRepository(db *database.Postgres, logger logger.Interface) (Repository, error) {
	if db == nil {
		return nil, errors.New("repository.NewUserRepository: db is null")
	}

	if logger == nil {
		return nil, errors.New("repository.NewUserRepository: logger is null")
	}

	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *repository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u == nil {
		return nil, errors.New("repository.Create: user is null")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	r.logger.Debug(fmt.Sprintf("txID: %s repository.Create, creating user: %v", txID, u))

	u.CreatedAt = time.Now()
	repoUsr := converter.ToUserFromDomain(u)

	query, args, err := sq.
		Insert("users").
		Columns("uuid", "login", "password", "created_at").
		Values(repoUsr.Uuid, repoUsr.Login, repoUsr.Password, repoUsr.CreatedAt).
		Suffix("RETURNING uuid").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s repository.Create, error building query: %v", txID, err))
		return nil, fmt.Errorf("repository.Create: %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&repoUsr.Uuid)
	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s repository.Create, error creating user: %v", txID, err))
		return nil, fmt.Errorf("repository.Create: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s repository.Create, successfully created user with UUID: %s", txID, repoUsr.Uuid))

	return converter.ToUserFromRepository(repoUsr), nil
}

func (r *repository) Delete(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s repository.Delete, deleting user with UUID: %s", txID, uuid))

	now := time.Now()
	query, args, err := sq.
		Update("users").
		Set("deleted_at", now).
		Where(sq.Eq{"uuid": uuid}).
		Suffix("RETURNING uuid, login, password, created_at, updated_at, deleted_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s repository.Delete, error building query: %v", txID, err))
		return nil, fmt.Errorf("repository.Delete: %w", err)
	}

	user := &repoUser.User{}
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.Uuid,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s repository.Delete, error soft deleting user: %v", txID, err))
		return nil, fmt.Errorf("repository.Delete: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s repository.Delete, successfully soft deleted user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}

func (r *repository) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s repository.Get, fetching user with UUID: %s", txID, uuid))

	query, args, err := sq.
		Select("uuid", "login", "password", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"uuid": uuid, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s repository.Get, error building query: %v", txID, err))
		return nil, fmt.Errorf("repository.Get: %w", err)
	}

	user := &repoUser.User{}
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.Uuid,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug(fmt.Sprintf("txID: %s repository.Get, user with UUID: %s not found", txID, uuid))
			return nil, nil
		}

		r.logger.Error(fmt.Sprintf("txID: %s repository.Get, error fetching user: %v", txID, err))
		return nil, fmt.Errorf("repository.Get: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s repository.Get, successfully fetched user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}
