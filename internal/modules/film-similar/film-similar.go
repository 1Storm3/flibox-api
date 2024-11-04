package film_similar

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-similar/handler"
	"kbox-api/internal/modules/film-similar/repository"
	"kbox-api/internal/modules/film-similar/service"
)

type Module struct {
	storage               *postgres.Storage
	config                *config.Config
	filmSimilarRepository service.FilmSimilarRepository
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
		f.filmSimilarService = service.NewFilmsSimilarService(repo, f.config, filmService)
	}
	return f.filmSimilarService, nil
}

func (f *Module) FilmSimilarRepository() (service.FilmSimilarRepository, error) {
	if f.filmSimilarRepository == nil {
		f.filmSimilarRepository = repository.NewFilmSimilarRepository(f.storage)
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
