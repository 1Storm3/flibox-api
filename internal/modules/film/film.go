package film

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/modules/film/handler"
	"kbox-api/internal/modules/film/repository"
	"kbox-api/internal/modules/film/service"
)

type ModuleInterface interface {
	FilmService() (service.FilmServiceInterface, error)
	FilmRepository() (repository.FilmRepositoryInterface, error)
	FilmHandler() (handler.FilmHandlerInterface, error)
}

type Module struct {
	storage        *postgres.Storage
	filmService    service.FilmServiceInterface
	filmRepository repository.FilmRepositoryInterface
	externalModule external.ModuleInterface
	filmHandler    handler.FilmHandlerInterface
}

func NewFilmModule(
	storage *postgres.Storage,
	externalModule external.ModuleInterface,
) ModuleInterface {
	return &Module{
		storage:        storage,
		externalModule: externalModule,
	}
}

func (f *Module) FilmService() (service.FilmServiceInterface, error) {
	if f.filmService == nil {
		repo, err := f.FilmRepository()

		if err != nil {
			return nil, err
		}

		externalService, err := f.externalModule.ExternalService()

		if err != nil {
			return nil, err
		}
		f.filmService = service.NewFilmService(repo, externalService)
	}

	return f.filmService, nil
}

func (f *Module) FilmRepository() (repository.FilmRepositoryInterface, error) {
	if f.filmRepository == nil {
		f.filmRepository = repository.NewFilmRepository(f.storage)
	}
	return f.filmRepository, nil

}

func (f *Module) FilmHandler() (handler.FilmHandlerInterface, error) {
	if f.filmHandler == nil {
		filmService, err := f.FilmService()
		if err != nil {
			return nil, err
		}
		f.filmHandler = handler.NewFilmHandler(filmService)
	}
	return f.filmHandler, nil
}
