package user

import (
	"app/internal/app/repository/postgres/user"
	"app/internal/domain"
	"app/internal/pkg/logger"
	"app/internal/pkg/types"
	"context"
	"errors"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	Delete(ctx context.Context, uuid string) error
	Get(ctx context.Context, uuid string) (*domain.User, error)
}

var _ Service = (*service)(nil)

type service struct {
	repository user.Repository
	logger     logger.Interface
}

func NewUserService(repository user.Repository, logger logger.Interface) (Service, error) {
	if repository == nil {
		return nil, errors.New("service.NewUserService: repository is null")
	}

	if logger == nil {
		return nil, errors.New("service.NewUserService: logger is null")
	}

	return &service{
		repository: repository,
		logger:     logger,
	}, nil
}

func (s *service) Create(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("service.Create: user is nil")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s service.Create, creating user: %v", txID, user))

	u, err := s.repository.Create(ctx, user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s service.Create, error creating user: %v", txID, err))
		return "", fmt.Errorf("service.Create: %w", err)
	}

	s.logger.Debug(fmt.Sprintf("txID: %s service.Create, successfully created user: %v", txID, user))
	return u.Uuid, nil
}

func (s *service) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s service.Get, fetching user with UUID: %s", txID, uuid))

	u, err := s.repository.Get(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s service.Get, error fetching user: %v", txID, err))
		return nil, fmt.Errorf("service.Get: %w", err)
	} else if u == nil {
		s.logger.Debug(fmt.Sprintf("txID: %s service.Get, user not found by UUID: %s", txID, uuid))
		return nil, domain.ErrorUserNotFound
	}

	s.logger.Debug(fmt.Sprintf("txID: %s service.Get, successfully fetched user: %v", txID, u))
	return u, nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s service.Delete, deleting user with UUID: %s", txID, uuid))

	u, err := s.repository.Delete(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s service.Delete, error deleting user: %v", txID, err))
		return fmt.Errorf("service.Delete: %w", err)
	}

	s.logger.Debug(fmt.Sprintf("txID: %s service.Delete, successfully deleted user: %v", txID, u))
	return nil
}
