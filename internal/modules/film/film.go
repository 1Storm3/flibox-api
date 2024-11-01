package film

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/external"
	"kinopoisk-api/internal/modules/film/handler"
	filmrepository "kinopoisk-api/internal/modules/film/repository"
	filmservice "kinopoisk-api/internal/modules/film/service"
)

type Module struct {
	storage        *postgres.Storage
	filmService    handler.FilmService
	filmRepository filmservice.FilmRepository
	externalModule *external.Module
	filmHandler    *handler.FilmHandler
}

func NewFilmModule(
	storage *postgres.Storage,
	externalModule *external.Module,
) *Module {
	return &Module{
		storage:        storage,
		externalModule: externalModule,
	}
}

func (f *Module) FilmService() (handler.FilmService, error) {
	if f.filmService == nil {
		repo, err := f.FilmRepository()

		if err != nil {
			return nil, err
		}

		externalService, err := f.externalModule.ExternalService()

		if err != nil {
			return nil, err
		}
		f.filmService = filmservice.NewFilmService(repo, externalService)
	}

	return f.filmService, nil
}

func (f *Module) FilmRepository() (filmservice.FilmRepository, error) {
	if f.filmRepository == nil {
		f.filmRepository = filmrepository.NewFilmRepository(f.storage)
	}
	return f.filmRepository, nil

}

func (f *Module) FilmHandler() (*handler.FilmHandler, error) {
	if f.filmHandler == nil {
		filmService, err := f.FilmService()
		if err != nil {
			return nil, err
		}
		f.filmHandler = handler.NewFilmHandler(filmService)
	}
	return f.filmHandler, nil
}
