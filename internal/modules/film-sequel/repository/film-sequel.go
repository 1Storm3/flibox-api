package repository

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/film-sequel/service"
	"kinopoisk-api/shared/httperror"
)

type filmSequelRepository struct {
	storage *postgres.Storage
}

func NewFilmSequelRepository(storage *postgres.Storage) *filmSequelRepository {
	return &filmSequelRepository{
		storage: storage,
	}
}

func (s *filmSequelRepository) GetAll(ctx context.Context, filmId string) ([]service.FilmSequel, error) {
	var filmSequels []service.FilmSequel
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmId).
		Preload("Film").
		Find(&filmSequels)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []service.FilmSequel{}, nil
	} else if result.Error != nil {
		return []service.FilmSequel{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return filmSequels, nil
}

func (s *filmSequelRepository) Save(filmId int, sequelId int) error {
	var existingSequel service.FilmSequel

	result := s.storage.DB().Where("film_id = ? AND sequel_id = ?", filmId, sequelId).First(&existingSequel)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().Create(&service.FilmSequel{
		FilmId:   filmId,
		SequelId: sequelId,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
