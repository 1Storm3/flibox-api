package user

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/modules/user/handler"
	"kbox-api/internal/modules/user/repository"
	"kbox-api/internal/modules/user/service"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	UserService() (service.UserServiceInterface, error)
	UserRepository() (repository.UserRepositoryInterface, error)
	UserHandler() (handler.UserHandlerInterface, error)
}

type Module struct {
	storage        *postgres.Storage
	userService    service.UserServiceInterface
	userRepository repository.UserRepositoryInterface
	userHandler    handler.UserHandlerInterface
	externalModule external.ModuleInterface
}

func NewUserModule(storage *postgres.Storage, externalModule external.ModuleInterface) ModuleInterface {
	return &Module{
		storage:        storage,
		externalModule: externalModule,
	}
}

func (u *Module) UserService() (service.UserServiceInterface, error) {
	if u.userService == nil {
		repo, err := u.UserRepository()
		if err != nil {
			return nil, err
		}
		s3Service, err := u.externalModule.S3Service()
		if err != nil {
			return nil, err
		}
		u.userService = service.NewUserService(repo, s3Service)
	}
	return u.userService, nil
}

func (u *Module) UserRepository() (repository.UserRepositoryInterface, error) {
	if u.userRepository == nil {
		u.userRepository = repository.NewUserRepository(u.storage)
	}
	return u.userRepository, nil
}

func (u *Module) UserHandler() (handler.UserHandlerInterface, error) {
	if u.userHandler == nil {
		userService, err := u.UserService()
		if err != nil {
			return nil, err
		}
		u.userHandler = handler.NewUserHandler(userService)
	}
	return u.userHandler, nil
}
