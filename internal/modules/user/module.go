package user

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/external"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	storage        *postgres.Storage
	service        ServiceInterface
	repository     RepositoryInterface
	handler        HandlerInterface
	externalModule external.ModuleInterface
}

func NewUserModule(storage *postgres.Storage, externalModule external.ModuleInterface) *Module {
	return &Module{
		storage:        storage,
		externalModule: externalModule,
	}
}

func (u *Module) Service() (ServiceInterface, error) {
	if u.service == nil {
		repo, err := u.Repository()
		if err != nil {
			return nil, err
		}
		s3Service, err := u.externalModule.S3Service()
		if err != nil {
			return nil, err
		}
		u.service = NewUserService(repo, s3Service)
	}
	return u.service, nil
}

func (u *Module) Repository() (RepositoryInterface, error) {
	if u.repository == nil {
		u.repository = NewUserRepository(u.storage)
	}
	return u.repository, nil
}

func (u *Module) Handler() (HandlerInterface, error) {
	if u.handler == nil {
		userService, err := u.Service()
		if err != nil {
			return nil, err
		}
		u.handler = NewUserHandler(userService)
	}
	return u.handler, nil
}
