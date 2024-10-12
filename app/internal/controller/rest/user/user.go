package rest

import (
	"app/domain"
	"context"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, uuid string) error
	Get(ctx context.Context, uuid string) (*domain.User, error)
}

type UserHandler struct {
	Service UserService
}
