package collectionfilm

import "github.com/1Storm3/flibox-api/database/postgres"

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Handler() (HandlerInterface, error)
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
}

type Module struct {
	storage    *postgres.Storage
	service    ServiceInterface
	repository RepositoryInterface
	handler    HandlerInterface
}

func NewCollectionFilmModule(storage *postgres.Storage) *Module {
	return &Module{
		storage: storage,
	}
}

func (m *Module) Handler() (HandlerInterface, error) {
	if m.handler == nil {
		collectionFilmService, err := m.Service()
		if err != nil {
			return nil, err
		}
		m.handler = NewCollectionFilmHandler(collectionFilmService)
	}
	return m.handler, nil
}

func (m *Module) Service() (ServiceInterface, error) {
	if m.service == nil {
		collectionFilmRepository, err := m.Repository()
		if err != nil {
			return nil, err
		}
		m.service = NewCollectionFilmService(collectionFilmRepository)
	}
	return m.service, nil
}

func (m *Module) Repository() (RepositoryInterface, error) {
	if m.repository == nil {
		m.repository = NewCollectionFilmRepository(m.storage)
	}
	return m.repository, nil
}
