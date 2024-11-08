package comment

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/modules/comment/handler"
	"kbox-api/internal/modules/comment/repository"
	"kbox-api/internal/modules/comment/service"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	CommentService() (service.CommentServiceInterface, error)
	CommentHandler() (handler.CommentHandlerInterface, error)
	CommentRepository() (repository.CommentRepositoryInterface, error)
}

type Module struct {
	storage           *postgres.Storage
	commentService    service.CommentServiceInterface
	commentHandler    handler.CommentHandlerInterface
	commentRepository repository.CommentRepositoryInterface
}

func NewCommentModule(storage *postgres.Storage) ModuleInterface {
	return &Module{
		storage: storage,
	}
}

func (c *Module) CommentService() (service.CommentServiceInterface, error) {
	if c.commentService == nil {
		repo, err := c.CommentRepository()
		if err != nil {
			return nil, err
		}
		c.commentService = service.NewCommentService(repo)
	}
	return c.commentService, nil
}
func (c *Module) CommentHandler() (handler.CommentHandlerInterface, error) {
	if c.commentHandler == nil {
		commentService, err := c.CommentService()
		if err != nil {
			return nil, err
		}
		c.commentHandler = handler.NewCommentHandler(commentService)
	}
	return c.commentHandler, nil
}
func (c *Module) CommentRepository() (repository.CommentRepositoryInterface, error) {
	if c.commentRepository == nil {
		c.commentRepository = repository.NewCommentRepository(c.storage)
	}
	return c.commentRepository, nil
}
