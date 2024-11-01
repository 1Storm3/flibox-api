package repository

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/film-similar/service"
	"kinopoisk-api/shared/httperror"
)

type filmSimilarRepository struct {
	storage *postgres.Storage
}

func NewFilmSimilarRepository(storage *postgres.Storage) *filmSimilarRepository {
	return &filmSimilarRepository{
		storage: storage,
	}
}

func (s *filmSimilarRepository) GetAll(ctx context.Context, filmId string) ([]service.FilmSimilar, error) {
	var filmSimilars []service.FilmSimilar

	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmId).
		Preload("Film").
		Find(&filmSimilars)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []service.FilmSimilar{}, nil
	} else if result.Error != nil {
		return []service.FilmSimilar{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error())
	}

	return filmSimilars, nil
}

func (s *filmSimilarRepository) Save(filmId int, similarId int) error {
	var existingSimilar service.FilmSimilar

	result := s.storage.DB().Where("film_id = ? AND similar_id = ?", filmId, similarId).First(&existingSimilar)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().Create(&service.FilmSimilar{
		FilmId:    filmId,
		SimilarId: similarId,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
