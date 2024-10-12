package converter

import (
	"app/domain"
	rest "app/internal/controller/rest/user/model"
)

func ToUserFromRest(user rest.User) (*domain.User, error) {
	return nil, nil
}

func ToUserFromDomain(user domain.User) (*rest.User, error) {
	return nil, nil
}
