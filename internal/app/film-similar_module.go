package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/modules/film-similar/handler"
	filmsimilarrepository "kinopoisk-api/internal/modules/film-similar/repository"
	filmsimilarservice "kinopoisk-api/internal/modules/film-similar/service"
)

type filmSimilarModule struct {
	storage               *postgres.Storage
	config                *config.Config
	filmSimilarRepository FilmSimilarRepository
	filmSimilarService    FilmSimilarService
	filmSimilarHandler    *handler.FilmSimilarHandler
	filmModule            *filmModule
}

func NewFilmSimilarModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule *filmModule,
) *filmSimilarModule {
	return &filmSimilarModule{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *filmSimilarModule) FilmSimilarService() (FilmSimilarService, error) {
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

func (f *filmSimilarModule) FilmSimilarRepository() (FilmSimilarRepository, error) {
	if f.filmSimilarRepository == nil {
		f.filmSimilarRepository = filmsimilarrepository.NewFilmSimilarRepository(f.storage)
	}
	return f.filmSimilarRepository, nil
}

func (f *filmSimilarModule) FilmSimilarHandler() (*handler.FilmSimilarHandler, error) {
	if f.filmSimilarHandler == nil {
		filmSimilarService, err := f.FilmSimilarService()
		if err != nil {
			return nil, err
		}
		f.filmSimilarHandler = handler.NewFilmSimilarHandler(filmSimilarService)
	}
	return f.filmSimilarHandler, nil
}
