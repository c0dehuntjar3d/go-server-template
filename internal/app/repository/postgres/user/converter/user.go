package converter

import (
	postgres "app/internal/app/repository/postgres/user/model"
	"app/internal/domain"
	"time"
)

func ToUserFromRepository(repoUser *postgres.User) *domain.User {
	var updatedAt *time.Time
	var deletedAt *time.Time

	if repoUser.UpdatedAt != nil && repoUser.UpdatedAt.Valid {
		updatedAt = &repoUser.UpdatedAt.Time
	}

	if repoUser.DeletedAt != nil && repoUser.DeletedAt.Valid {
		deletedAt = &repoUser.DeletedAt.Time
	}

	return &domain.User{
		Uuid:      repoUser.Uuid,
		Login:     repoUser.Login,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func ToUserFromDomain(domainUser *domain.User) *postgres.User {
	return &postgres.User{
		Uuid:      domainUser.Uuid,
		Login:     domainUser.Login,
		Password:  domainUser.Password,
		CreatedAt: domainUser.CreatedAt,
	}
}
