package postgres

import (
	"app/domain"
	"context"
)

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) Create(_ context.Context, userUUID string, info *domain.UserInfo) error {

	return nil
}

func (r *repository) Get(_ context.Context, uuid string) (*domain.User, error) {
	return nil, nil
}
