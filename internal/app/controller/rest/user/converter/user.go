package converter

import (
	"app/internal/app/controller/rest/user/model"
	"app/internal/domain"
	"time"

	"github.com/google/uuid"
)

func ToUserFromRest(user rest.User) *domain.User {
	return &domain.User{
		Uuid:      uuid.NewString(),
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}
}
