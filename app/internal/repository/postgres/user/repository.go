package user

import (
	"app/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, uuid string) (*domain.User, error)
	Delete(ctx context.Context, uuid string) (*domain.User, error)
}
