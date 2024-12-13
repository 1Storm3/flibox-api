package filmsequel

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/internal/shared/httperror"
)

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error)
	Save(ctx context.Context, filmId int, sequelId int) error
}

type Repository struct {
	storage *postgres.Storage
}

func NewFilmSequelRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (s *Repository) GetAll(ctx context.Context, filmID string) ([]model.FilmSequel, error) {
	var filmSequels []model.FilmSequel
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film").
		Find(&filmSequels)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.FilmSequel{}, nil
	} else if result.Error != nil {
		return []model.FilmSequel{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return filmSequels, nil
}

func (s *Repository) Save(ctx context.Context, filmID int, sequelID int) error {
	var existingSequel model.FilmSequel

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND sequel_id = ?", filmID, sequelID).First(&existingSequel)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().Create(&model.FilmSequel{
		FilmID:   filmID,
		SequelID: sequelID,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
