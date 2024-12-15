package filmsequel

import (
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/modules/film"
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
	filmModule film.ModuleInterface
	handler    HandlerInterface
}

func NewFilmSequelModule(
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
		f.service = NewFilmsSequelService(repo, f.config, filmService)
	}
	return f.service, nil
}

func (f *Module) Repository() (RepositoryInterface, error) {
	if f.repository == nil {
		f.repository = NewFilmSequelRepository(f.storage)
	}
	return f.repository, nil
}

func (f *Module) Handler() (HandlerInterface, error) {
	if f.handler == nil {
		filmSequelService, err := f.Service()
		if err != nil {
			return nil, err
		}
		f.handler = NewFilmSequelHandler(filmSequelService)
	}
	return f.handler, nil
}
