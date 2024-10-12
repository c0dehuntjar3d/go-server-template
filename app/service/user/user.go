package user

import (
	"app/domain"
	def "app/internal/controller/rest/user"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, uuid string) (*domain.User, error)
	Delete(ctx context.Context, uuid string) (*domain.User, error)
}

var _ def.UserService = (*service)(nil)

type service struct {
}

func (s *service) Create(ctx context.Context, user *domain.User) error {
	panic("unimplemented")
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	panic("unimplemented")
}

func (s *service) Get(ctx context.Context, uuid string) (*domain.User, error) {
	panic("unimplemented")
}
