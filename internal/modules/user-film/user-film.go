package user_film

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/user-film/handler"
	"kbox-api/internal/modules/user-film/repository"
	"kbox-api/internal/modules/user-film/service"
)

type Module struct {
	storage            *postgres.Storage
	userFilmRepository service.UserFilmRepository
	userFilmService    handler.UserFilmService
	userFilmHandler    *handler.UserFilmHandler
	filmModule         *film.Module
}

func NewUserFilmModule(
	storage *postgres.Storage,
	filmModule *film.Module,
) *Module {
	return &Module{
		storage:    storage,
		filmModule: filmModule,
	}
}

func (u *Module) UserFilmService() (handler.UserFilmService, error) {
	if u.userFilmService == nil {
		repo, err := u.UserFilmRepository()
		if err != nil {
			return nil, err
		}
		u.userFilmService = service.NewUserFilmService(repo)
	}
	return u.userFilmService, nil
}

func (u *Module) UserFilmRepository() (service.UserFilmRepository, error) {
	if u.userFilmRepository == nil {
		u.userFilmRepository = repository.NewUserFilmRepository(u.storage)
	}
	return u.userFilmRepository, nil
}

func (u *Module) UserFilmHandler() (*handler.UserFilmHandler, error) {
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
