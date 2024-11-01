package service

import "context"

type UserRepository interface {
	GetOneById(ctx context.Context, id string) (User, error)
	GetOneByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)
	GetOneByNickName(ctx context.Context, nickName string) (User, error)
}
