package service

import (
	"context"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/user/dto"
)

type UserRepository interface {
	GetOneById(ctx context.Context, id string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error)
}
