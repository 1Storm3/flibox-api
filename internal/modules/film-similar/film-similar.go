package film_similar

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/modules/film"
	"kinopoisk-api/internal/modules/film-similar/handler"
	filmsimilarrepository "kinopoisk-api/internal/modules/film-similar/repository"
	filmsimilarservice "kinopoisk-api/internal/modules/film-similar/service"
)

type Module struct {
	storage               *postgres.Storage
	config                *config.Config
	filmSimilarRepository filmsimilarservice.FilmSimilarRepository
	filmSimilarService    handler.FilmSimilarService
	filmSimilarHandler    *handler.FilmSimilarHandler
	filmModule            *film.Module
}

func NewFilmSimilarModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule *film.Module,
) *Module {
	return &Module{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *Module) FilmSimilarService() (handler.FilmSimilarService, error) {
	if f.filmSimilarService == nil {
		repo, err := f.FilmSimilarRepository()
		if err != nil {
			return nil, err
		}
		filmService, err := f.filmModule.FilmService()
		if err != nil {
			return nil, err
		}
		f.filmSimilarService = filmsimilarservice.NewFilmsSimilarService(repo, f.config, filmService)
	}
	return f.filmSimilarService, nil
}

func (f *Module) FilmSimilarRepository() (filmsimilarservice.FilmSimilarRepository, error) {
	if f.filmSimilarRepository == nil {
		f.filmSimilarRepository = filmsimilarrepository.NewFilmSimilarRepository(f.storage)
	}
	return f.filmSimilarRepository, nil
}

func (f *Module) FilmSimilarHandler() (*handler.FilmSimilarHandler, error) {
	if f.filmSimilarHandler == nil {
		filmSimilarService, err := f.FilmSimilarService()
		if err != nil {
			return nil, err
		}
		f.filmSimilarHandler = handler.NewFilmSimilarHandler(filmSimilarService)
	}
	return f.filmSimilarHandler, nil
}
