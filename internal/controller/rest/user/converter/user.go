package converter

import (
	"time"

	"github.com/google/uuid"

	"app/domain"
	rest "app/internal/controller/rest/user/model"
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
