package collection

import (
	"kbox-api/database/postgres"
)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	storage    *postgres.Storage
	service    ServiceInterface
	repository RepositoryInterface
	handler    HandlerInterface
}

func NewCollectionModule(storage *postgres.Storage) *Module {
	return &Module{
		storage: storage,
	}
}

func (m *Module) Service() (ServiceInterface, error) {
	if m.service == nil {
		repo, err := m.Repository()
		if err != nil {
			return nil, err
		}
		m.service = NewCollectionService(repo)
	}
	return m.service, nil
}

func (m *Module) Repository() (RepositoryInterface, error) {
	if m.repository == nil {
		m.repository = NewCollectionRepository(m.storage)
	}
	return m.repository, nil
}

func (m *Module) Handler() (HandlerInterface, error) {
	if m.handler == nil {
		collectionService, err := m.Service()
		if err != nil {
			return nil, err
		}
		m.handler = NewCollectionHandler(collectionService)
	}
	return m.handler, nil
}
