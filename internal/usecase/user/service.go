package user

import (
	"context"

	"app/domain"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	Delete(ctx context.Context, uuid string) error
	Get(ctx context.Context, uuid string) (*domain.User, error)
}
