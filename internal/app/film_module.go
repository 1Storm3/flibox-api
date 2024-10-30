package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/film/handler"
	filmrepository "kinopoisk-api/internal/modules/film/repository"
	filmservice "kinopoisk-api/internal/modules/film/service"
)

type filmModule struct {
	storage         *postgres.Storage
	filmService     FilmService
	filmRepository  FilmRepository
	externalService ExternalService
	filmHandler     *handler.FilmHandler
}

func NewFilmModule(
	storage *postgres.Storage,
	externalService ExternalService,
) *filmModule {
	return &filmModule{
		storage:         storage,
		externalService: externalService,
	}
}

func (f *filmModule) FilmService() (FilmService, error) {
	if f.filmService == nil {
		repo, err := f.FilmRepository()

		if err != nil {
			return nil, err
		}
		f.filmService = filmservice.NewFilmService(repo, f.externalService)
	}

	return f.filmService, nil
}

func (f *filmModule) FilmRepository() (FilmRepository, error) {
	if f.filmRepository == nil {
		f.filmRepository = filmrepository.NewFilmRepository(f.storage)
	}
	return f.filmRepository, nil

}

func (f *filmModule) FilmHandler() (*handler.FilmHandler, error) {
	if f.filmHandler == nil {
		filmService, err := f.FilmService()
		if err != nil {
			return nil, err
		}
		f.filmHandler = handler.NewFilmHandler(filmService)
	}
	return f.filmHandler, nil
}
