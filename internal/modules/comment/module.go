package comment

import (
	"kbox-api/database/postgres"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Handler() (HandlerInterface, error)
	Repository() (RepositoryInterface, error)
}

type Module struct {
	storage    *postgres.Storage
	service    ServiceInterface
	handler    HandlerInterface
	repository RepositoryInterface
}

func NewCommentModule(storage *postgres.Storage) *Module {
	return &Module{
		storage: storage,
	}
}

func (c *Module) Service() (ServiceInterface, error) {
	if c.service == nil {
		repo, err := c.Repository()
		if err != nil {
			return nil, err
		}
		c.service = NewCommentService(repo)
	}
	return c.service, nil
}
func (c *Module) Handler() (HandlerInterface, error) {
	if c.handler == nil {
		commentService, err := c.Service()
		if err != nil {
			return nil, err
		}
		c.handler = NewCommentHandler(commentService)
	}
	return c.handler, nil
}
func (c *Module) Repository() (RepositoryInterface, error) {
	if c.repository == nil {
		c.repository = NewCommentRepository(c.storage)
	}
	return c.repository, nil
}
