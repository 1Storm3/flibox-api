package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user-film/handler"
	userfilmrepository "kinopoisk-api/internal/modules/user-film/repository"
	userfilmservice "kinopoisk-api/internal/modules/user-film/service"
)

type userFilmModule struct {
	storage            *postgres.Storage
	userFilmRepository UserFilmRepository
	userFilmService    UserFilmService
	userFilmHandler    *handler.UserFilmHandler
	filmModule         *filmModule
}

func NewUserFilmModule(
	storage *postgres.Storage,
	filmModule *filmModule,
) *userFilmModule {
	return &userFilmModule{
		storage:    storage,
		filmModule: filmModule,
	}
}

func (u *userFilmModule) UserFilmService() (UserFilmService, error) {
	if u.userFilmService == nil {
		repo, err := u.UserFilmRepository()
		if err != nil {
			return nil, err
		}
		u.userFilmService = userfilmservice.NewUserFilmService(repo)
	}
	return u.userFilmService, nil
}

func (u *userFilmModule) UserFilmRepository() (UserFilmRepository, error) {
	if u.userFilmRepository == nil {
		u.userFilmRepository = userfilmrepository.NewUserFilmRepository(u.storage)
	}
	return u.userFilmRepository, nil
}

func (u *userFilmModule) UserFilmHandler() (*handler.UserFilmHandler, error) {
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
