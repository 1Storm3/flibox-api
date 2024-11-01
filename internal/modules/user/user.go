package user

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user/handler"
	userrepository "kinopoisk-api/internal/modules/user/repository"
	userservice "kinopoisk-api/internal/modules/user/service"
)

type Module struct {
	storage        *postgres.Storage
	userService    handler.UserService
	userRepository userservice.UserRepository
	userHandler    *handler.UserHandler
}

func NewUserModule(storage *postgres.Storage) *Module {
	return &Module{
		storage: storage,
	}
}

func (u *Module) UserService() (handler.UserService, error) {
	if u.userService == nil {
		repo, err := u.UserRepository()
		if err != nil {
			return nil, err
		}
		u.userService = userservice.NewUserService(repo)
	}
	return u.userService, nil
}

func (u *Module) UserRepository() (userservice.UserRepository, error) {
	if u.userRepository == nil {
		u.userRepository = userrepository.NewUserRepository(u.storage)
	}
	return u.userRepository, nil
}

func (u *Module) UserHandler() (*handler.UserHandler, error) {
	if u.userHandler == nil {
		userService, err := u.UserService()
		if err != nil {
			return nil, err
		}
		u.userHandler = handler.NewUserHandler(userService)
	}
	return u.userHandler, nil
}
