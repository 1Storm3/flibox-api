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

type FilmSequelRepositoryInterface interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error)
	Save(filmId int, sequelId int) error
}

type filmSequelRepository struct {
	storage *postgres.Storage
}

func NewFilmSequelRepository(storage *postgres.Storage) FilmSequelRepositoryInterface {
	return &filmSequelRepository{
		storage: storage,
	}
}

func (s *filmSequelRepository) GetAll(ctx context.Context, filmID string) ([]model.FilmSequel, error) {
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

func (s *filmSequelRepository) Save(filmID int, sequelID int) error {
	var existingSequel model.FilmSequel

	result := s.storage.DB().Where("film_id = ? AND sequel_id = ?", filmID, sequelID).First(&existingSequel)

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
