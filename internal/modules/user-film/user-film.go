package userfilm

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/user-film/handler"
	"kbox-api/internal/modules/user-film/repository"
	"kbox-api/internal/modules/user-film/service"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	UserFilmService() (service.UserFilmServiceInterface, error)
	UserFilmRepository() (repository.UserFilmRepositoryInterface, error)
	UserFilmHandler() (handler.UserFilmHandlerInterface, error)
}

type Module struct {
	storage            *postgres.Storage
	userFilmRepository repository.UserFilmRepositoryInterface
	userFilmService    service.UserFilmServiceInterface
	userFilmHandler    handler.UserFilmHandlerInterface
	filmModule         film.ModuleInterface
}

func NewUserFilmModule(
	storage *postgres.Storage,
	filmModule film.ModuleInterface,
) *Module {
	return &Module{
		storage:    storage,
		filmModule: filmModule,
	}
}

func (u *Module) UserFilmService() (service.UserFilmServiceInterface, error) {
	if u.userFilmService == nil {
		repo, err := u.UserFilmRepository()
		if err != nil {
			return nil, err
		}
		u.userFilmService = service.NewUserFilmService(repo)
	}
	return u.userFilmService, nil
}

func (u *Module) UserFilmRepository() (repository.UserFilmRepositoryInterface, error) {
	if u.userFilmRepository == nil {
		u.userFilmRepository = repository.NewUserFilmRepository(u.storage)
	}
	return u.userFilmRepository, nil
}

func (u *Module) UserFilmHandler() (handler.UserFilmHandlerInterface, error) {
	if u.userFilmHandler == nil {
		userFilmService, err := u.UserFilmService()
		if err != nil {
			return nil, err
		}
		filmService, err := u.filmModule.FilmService()
		if err != nil {
			return nil, err
		}
		u.userFilmHandler = handler.NewUserFilmHandler(userFilmService, filmService)
	}
	return u.userFilmHandler, nil
}
