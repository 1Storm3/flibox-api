package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user/handler"
	userrepository "kinopoisk-api/internal/modules/user/repository"
	userservice "kinopoisk-api/internal/modules/user/service"
)

type userModule struct {
	storage        *postgres.Storage
	userService    UserService
	userRepository UserRepository
	userHandler    *handler.UserHandler
}

func NewUserModule(storage *postgres.Storage) *userModule {
	return &userModule{
		storage: storage,
	}
}

func (u *userModule) UserService() (UserService, error) {
	if u.userService == nil {
		repo, err := u.UserRepository()
		if err != nil {
			return nil, err
		}
		u.userService = userservice.NewUserService(repo)
	}
	return u.userService, nil
}

func (u *userModule) UserRepository() (UserRepository, error) {
	if u.userRepository == nil {
		u.userRepository = userrepository.NewUserRepository(u.storage)
	}
	return u.userRepository, nil
}

func (u *userModule) UserHandler() (*handler.UserHandler, error) {
	if u.userHandler == nil {
		userService, err := u.UserService()
		if err != nil {
			return nil, err
		}
		u.userHandler = handler.NewUserHandler(userService)
	}
	return u.userHandler, nil
}
