package user

import (
	"app/domain"
	def "app/internal/controller/rest/user"
	"app/pkg/logger"
	"app/pkg/types"
	"context"
	"errors"
	"fmt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, uuid string) (*domain.User, error)
	Delete(ctx context.Context, uuid string) (*domain.User, error)
}

var _ def.UserService = (*service)(nil)

type service struct {
	repository UserRepository
	logger     logger.Interface
}

func NewUserService(repository UserRepository, logger logger.Interface) (*service, error) {
	if repository == nil {
		return nil, errors.New("db is null")
	}

	if logger == nil {
		return nil, errors.New("logger is null")
	}

	return &service{
		repository: repository,
		logger:     logger,
	}, nil
}

func (s *service) Create(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s, Creating user: %v", txID, user))

	u, err := s.repository.Create(ctx, user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s, Error creating user: %v", txID, err))
		return "", err
	}

	s.logger.Debug(fmt.Sprintf("txID: %s, Successfully created user: %v", txID, user))
	return u.Uuid, nil
}

func (s *service) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s, Fetching user with UUID: %s", txID, uuid))

	user, err := s.repository.Get(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s, Error fetching user: %v", txID, err))
		return nil, err
	}

	s.logger.Debug(fmt.Sprintf("txID: %s, Successfully fetched user: %v", txID, user))
	return user, nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s, Soft deleting user with UUID: %s", txID, uuid))

	user, err := s.repository.Delete(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s, Error soft deleting user: %v", txID, err))
		return err
	}

	s.logger.Debug(fmt.Sprintf("txID: %s, Successfully soft deleted user: %v", txID, user))
	return nil
}
