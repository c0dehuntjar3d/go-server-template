package converter

import (
	"app/domain"
	rest "app/internal/controller/rest/user/model"
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
