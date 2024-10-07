package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
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
		return []service.FilmSimilar{}, result.Error
	}

	return filmSimilars, nil
}

func (s *filmSimilarRepository) Save(filmId int, similarId int) error {
	var existingSimilar service.FilmSimilar

	result := s.storage.DB().Where("film_id = ? AND similar_id = ?", filmId, similarId).First(&existingSimilar)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	createResult := s.storage.DB().Create(&service.FilmSimilar{
		FilmId:    filmId,
		SimilarId: similarId,
	})

	if createResult.Error != nil {
		return createResult.Error
	}

	return nil
}
