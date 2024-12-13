package filmsimilar

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	storage    *postgres.Storage
	config     *config.Config
	repository RepositoryInterface
	service    ServiceInterface
	handler    HandlerInterface
	filmModule film.ModuleInterface
}

func NewFilmSimilarModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule film.ModuleInterface,
) *Module {
	return &Module{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *Module) Service() (ServiceInterface, error) {
	if f.service == nil {
		repo, err := f.Repository()
		if err != nil {
			return nil, err
		}
		filmService, err := f.filmModule.Service()
		if err != nil {
			return nil, err
		}
		f.service = NewFilmsSimilarService(repo, f.config, filmService)
	}
	return f.service, nil
}

func (f *Module) Repository() (RepositoryInterface, error) {
	if f.repository == nil {
		f.repository = NewFilmSimilarRepository(f.storage)
	}
	return f.repository, nil
}

func (f *Module) Handler() (HandlerInterface, error) {
	if f.handler == nil {
		filmSimilarService, err := f.Service()
		if err != nil {
			return nil, err
		}
		f.handler = NewFilmSimilarHandler(filmSimilarService)
	}
	return f.handler, nil
}
