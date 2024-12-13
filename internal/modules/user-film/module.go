package userfilm

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/recommendation/adapter"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	storage                *postgres.Storage
	repository             RepositoryInterface
	service                ServiceInterface
	handler                HandlerInterface
	filmModule             film.ModuleInterface
	recommendModuleFactory func() (adapter.ModuleInterface, error)
}

func NewUserFilmModule(
	storage *postgres.Storage,
	filmModule film.ModuleInterface,
	recommendModuleFactory func() (adapter.ModuleInterface, error),
) *Module {
	return &Module{
		storage:                storage,
		filmModule:             filmModule,
		recommendModuleFactory: recommendModuleFactory,
	}
}

func (u *Module) Service() (ServiceInterface, error) {
	if u.service == nil {
		repo, err := u.Repository()
		if err != nil {
			return nil, err
		}
		u.service = NewUserFilmService(repo)
	}
	return u.service, nil
}

func (u *Module) Repository() (RepositoryInterface, error) {
	if u.repository == nil {
		u.repository = NewUserFilmRepository(u.storage)
	}
	return u.repository, nil
}

func (u *Module) Handler() (HandlerInterface, error) {
	if u.handler == nil {
		userFilmService, err := u.Service()
		if err != nil {
			return nil, err
		}

		filmService, err := u.filmModule.Service()
		if err != nil {
			return nil, err
		}
		recommendModule, err := u.recommendModuleFactory()
		if err != nil {
			return nil, err
		}
		recommendService, err := recommendModule.Service()
		if err != nil {
			return nil, err
		}

		u.handler = NewUserFilmHandler(userFilmService, filmService, recommendService)
	}
	return u.handler, nil
}
