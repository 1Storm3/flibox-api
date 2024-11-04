package handler

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/user/dto"
)

type UserService interface {
	GetOneByNickName(nickName string) (model.User, error)
	GetOneByEmail(email string) (model.User, error)
	CheckPassword(user model.User, password string) bool
	HashPassword(password string) (string, error)
	Create(user model.User) (model.User, error)
	GetOneById(id string) (model.User, error)
	Update(userDTO dto.UpdateUserDTO) (model.User, error)
}
