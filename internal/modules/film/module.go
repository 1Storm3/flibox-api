package film

import (
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/modules/external"
)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	storage        *postgres.Storage
	service        ServiceInterface
	repository     RepositoryInterface
	externalModule external.ModuleInterface
	handler        HandlerInterface
}

func NewFilmModule(
	storage *postgres.Storage,
	externalModule external.ModuleInterface,
) *Module {
	return &Module{
		storage:        storage,
		externalModule: externalModule,
	}
}

func (f *Module) Service() (ServiceInterface, error) {
	if f.service == nil {
		repo, err := f.Repository()

		if err != nil {
			return nil, err
		}

		externalService, err := f.externalModule.Service()

		if err != nil {
			return nil, err
		}
		f.service = NewFilmService(repo, externalService)
	}

	return f.service, nil
}

func (f *Module) Repository() (RepositoryInterface, error) {
	if f.repository == nil {
		f.repository = NewFilmRepository(f.storage)
	}
	return f.repository, nil

}

func (f *Module) Handler() (HandlerInterface, error) {
	if f.handler == nil {
		filmService, err := f.Service()
		if err != nil {
			return nil, err
		}

		f.handler = NewFilmHandler(filmService)
	}
	return f.handler, nil
}
