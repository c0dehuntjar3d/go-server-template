package user

import (
	"context"
	"errors"
	"fmt"

	"app/domain"
	"app/internal/repository/postgres/user"
	"app/pkg/logger"
	"app/pkg/types"
)

var _ UserService = (*userService)(nil)

type userService struct {
	repository user.UserRepository
	logger     logger.Interface
}

func NewUserService(repository user.UserRepository, logger logger.Interface) (*userService, error) {
	if repository == nil {
		return nil, errors.New("userService.NewUserService: repository is null")
	}

	if logger == nil {
		return nil, errors.New("userService.NewUserService: logger is null")
	}

	return &userService{
		repository: repository,
		logger:     logger,
	}, nil
}

func (s *userService) Create(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("userService.Create: user is nil")
	}

	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s userService.Create, creating user: %v", txID, user))

	u, err := s.repository.Create(ctx, user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s userService.Create, error creating user: %v", txID, err))
		return "", fmt.Errorf("userService.Create: %w", err)
	}

	s.logger.Debug(fmt.Sprintf("txID: %s userService.Create, successfully created user: %v", txID, user))
	return u.Uuid, nil
}

func (s *userService) Get(ctx context.Context, uuid string) (*domain.User, error) {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s userService.Get, fetching user with UUID: %s", txID, uuid))

	user, err := s.repository.Get(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s userService.Get, error fetching user: %v", txID, err))
		return nil, fmt.Errorf("userService.Get: %w", err)
	} else if user == nil {
		s.logger.Debug(fmt.Sprintf("txID: %s userService.Get, user not found by UUID: %s", txID, uuid))
		return nil, domain.ErrorUserNotFound
	}

	s.logger.Debug(fmt.Sprintf("txID: %s userService.Get, successfully fetched user: %v", txID, user))
	return user, nil
}

func (s *userService) Delete(ctx context.Context, uuid string) error {
	txID := ctx.Value(types.CtxKey("tx")).(string)
	s.logger.Debug(fmt.Sprintf("txID: %s userService.Delete, deleting user with UUID: %s", txID, uuid))

	user, err := s.repository.Delete(ctx, uuid)
	if err != nil {
		s.logger.Error(fmt.Sprintf("txID: %s userService.Delete, error deleting user: %v", txID, err))
		return fmt.Errorf("userService.Delete: %w", err)
	}

	s.logger.Debug(fmt.Sprintf("txID: %s userService.Delete, successfully deleted user: %v", txID, user))
	return nil
}
