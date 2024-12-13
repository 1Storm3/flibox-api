package filmsimilar

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
	GetAll(ctx context.Context, filmId string) ([]model.FilmSimilar, error)
	Save(ctx context.Context, filmId int, similarId int) error
}

type Repository struct {
	storage *postgres.Storage
}

func NewFilmSimilarRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (s *Repository) GetAll(ctx context.Context, filmID string) ([]model.FilmSimilar, error) {
	var filmSimilars []model.FilmSimilar

	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film").
		Find(&filmSimilars)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.FilmSimilar{}, nil
	} else if result.Error != nil {
		return []model.FilmSimilar{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error())
	}

	return filmSimilars, nil
}

func (s *Repository) Save(ctx context.Context, filmID int, similarID int) error {
	var existingSimilar model.FilmSimilar

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND similar_id = ?", filmID, similarID).First(&existingSimilar)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().WithContext(ctx).Create(&model.FilmSimilar{
		FilmID:    filmID,
		SimilarID: similarID,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
