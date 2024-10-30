package handler

import userservice "kinopoisk-api/internal/modules/user/service"

type UserService interface {
	GetOne(userToken string) (userservice.User, error)
}
