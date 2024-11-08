package filmsimilar

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-similar/handler"
	"kbox-api/internal/modules/film-similar/repository"
	"kbox-api/internal/modules/film-similar/service"
)

type ModuleInterface interface {
	FilmSimilarService() (service.FilmSimilarServiceInterface, error)
	FilmSimilarRepository() (repository.FilmSimilarRepositoryInterface, error)
	FilmSimilarHandler() (handler.FilmSimilarHandlerInterface, error)
}

type Module struct {
	storage               *postgres.Storage
	config                *config.Config
	filmSimilarRepository repository.FilmSimilarRepositoryInterface
	filmSimilarService    service.FilmSimilarServiceInterface
	filmSimilarHandler    handler.FilmSimilarHandlerInterface
	filmModule            film.ModuleInterface
}

func NewFilmSimilarModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule film.ModuleInterface,
) ModuleInterface {
	return &Module{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *Module) FilmSimilarService() (service.FilmSimilarServiceInterface, error) {
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

func (f *Module) FilmSimilarRepository() (repository.FilmSimilarRepositoryInterface, error) {
	if f.filmSimilarRepository == nil {
		f.filmSimilarRepository = repository.NewFilmSimilarRepository(f.storage)
	}
	return f.filmSimilarRepository, nil
}

func (f *Module) FilmSimilarHandler() (handler.FilmSimilarHandlerInterface, error) {
	if f.filmSimilarHandler == nil {
		filmSimilarService, err := f.FilmSimilarService()
		if err != nil {
			return nil, err
		}
		f.filmSimilarHandler = handler.NewFilmSimilarHandler(filmSimilarService)
	}
	return f.filmSimilarHandler, nil
}
