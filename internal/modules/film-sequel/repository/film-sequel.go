package repository

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/shared/httperror"
)

type filmSequelRepository struct {
	storage *postgres.Storage
}

func NewFilmSequelRepository(storage *postgres.Storage) *filmSequelRepository {
	return &filmSequelRepository{
		storage: storage,
	}
}

func (s *filmSequelRepository) GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error) {
	var filmSequels []model.FilmSequel
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmId).
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

func (s *filmSequelRepository) Save(filmId int, sequelId int) error {
	var existingSequel model.FilmSequel

	result := s.storage.DB().Where("film_id = ? AND sequel_id = ?", filmId, sequelId).First(&existingSequel)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().Create(&model.FilmSequel{
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
