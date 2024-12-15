package historyfilms

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/model"
)

type ServiceInterface interface {
	Add(ctx context.Context, filmId, userId string) error
	GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewHistoryFilmsService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Add(ctx context.Context, filmId, userId string) error {
	return s.repository.Add(ctx, filmId, userId)
}

func (s *Service) GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error) {
	return s.repository.GetAll(ctx, userId)
}
