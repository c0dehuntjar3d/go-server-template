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

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func NewRepository() *userRepository {
	return &userRepository{}
}

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	db     *database.Postgres
	logger logger.Interface
}

func NewUserRepository(db *database.Postgres, logger logger.Interface) (*userRepository, error) {
	if db == nil {
		return nil, errors.New("userRepository.NewUserRepository: db is null")
	}

	if logger == nil {
		return nil, errors.New("userRepository.NewUserRepository: logger is null")
	}

	return &userRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u == nil {
		return nil, errors.New("userRepository.Create: user is null")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Create, creating user: %v", txID, u))

	u.CreatedAt = time.Now()
	repoUser := converter.ToUserFromDomain(u)

	query, args, err := sq.
		Insert("users").
		Columns("uuid", "login", "password", "created_at").
		Values(repoUser.Uuid, repoUser.Login, repoUser.Password, repoUser.CreatedAt).
		Suffix("RETURNING uuid").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Create, error building query: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Create: %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&repoUser.Uuid)
	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Create, error creating user: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Create: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Create, successfully created user with UUID: %s", txID, repoUser.Uuid))

	return converter.ToUserFromRepository(repoUser), nil
}

func (r *userRepository) Delete(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Delete, deleting user with UUID: %s", txID, uuid))

	now := time.Now()
	query, args, err := sq.
		Update("users").
		Set("deleted_at", now).
		Where(sq.Eq{"uuid": uuid}).
		Suffix("RETURNING uuid, login, password, created_at, updated_at, deleted_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Delete, error building query: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Delete: %w", err)
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
		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Delete, error soft deleting user: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Delete: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Delete, successfully soft deleted user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}

func (r *userRepository) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)

	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Get, fetching user with UUID: %s", txID, uuid))

	query, args, err := sq.
		Select("uuid", "login", "password", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"uuid": uuid, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Get, error building query: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Get: %w", err)
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
			r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Get, user with UUID: %s not found", txID, uuid))
			return nil, nil
		}

		r.logger.Error(fmt.Sprintf("txID: %s userRepository.Get, error fetching user: %v", txID, err))
		return nil, fmt.Errorf("userRepository.Get: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("txID: %s userRepository.Get, successfully fetched user with UUID: %s", txID, user.Uuid))

	return converter.ToUserFromRepository(user), nil
}
