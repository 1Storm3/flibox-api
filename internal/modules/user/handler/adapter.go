package handler

import userservice "kinopoisk-api/internal/modules/user/service"

type UserService interface {
	GetOneByNickName(nickName string) (userservice.User, error)
	GetOneByEmail(email string) (userservice.User, error)
	CheckPassword(user userservice.User, password string) bool
	HashPassword(password string) (string, error)
	CreateUser(user userservice.User) (userservice.User, error)
	GetOneById(id string) (userservice.User, error)
}
