package historyfilms

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/recommendation/adapter"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Repository() (RepositoryInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	repository RepositoryInterface
	storage    *postgres.Storage
	service    ServiceInterface
	handler    HandlerInterface

	recommendModuleFactory func() (adapter.ModuleInterface, error)
}

func NewHistoryFilmsModule(
	storage *postgres.Storage,
	recommendModuleFactory func() (adapter.ModuleInterface, error),
) *Module {
	return &Module{
		storage:                storage,
		recommendModuleFactory: recommendModuleFactory,
	}
}

func (h *Module) Service() (ServiceInterface, error) {
	if h.service == nil {
		repo, err := h.Repository()
		if err != nil {
			return nil, err
		}
		h.service = NewHistoryFilmsService(repo)
	}
	return h.service, nil
}

func (h *Module) Repository() (RepositoryInterface, error) {
	if h.repository == nil {
		h.repository = NewHistoryFilmsRepository(h.storage)
	}
	return h.repository, nil
}

func (h *Module) Handler() (HandlerInterface, error) {
	if h.handler == nil {
		service, err := h.Service()
		if err != nil {
			return nil, err
		}
		recommendModule, err := h.recommendModuleFactory()
		if err != nil {
			return nil, err
		}
		recommendService, err := recommendModule.Service()
		if err != nil {
			return nil, err
		}
		h.handler = NewHistoryFilmsHandler(service, recommendService)
	}
	return h.handler, nil
}
