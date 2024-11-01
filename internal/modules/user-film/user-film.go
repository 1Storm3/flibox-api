package user_film

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/film"
	"kinopoisk-api/internal/modules/user-film/handler"
	userfilmrepository "kinopoisk-api/internal/modules/user-film/repository"
	userfilmservice "kinopoisk-api/internal/modules/user-film/service"
)

type Module struct {
	storage            *postgres.Storage
	userFilmRepository userfilmservice.UserFilmRepository
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
		u.userFilmService = userfilmservice.NewUserFilmService(repo)
	}
	return u.userFilmService, nil
}

func (u *Module) UserFilmRepository() (userfilmservice.UserFilmRepository, error) {
	if u.userFilmRepository == nil {
		u.userFilmRepository = userfilmrepository.NewUserFilmRepository(u.storage)
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
