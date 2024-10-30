package service

import "context"

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (User, error)
}
