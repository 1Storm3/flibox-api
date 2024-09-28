package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
)

type filmSequelRepository struct {
	storage *postgres.Storage
}

func NewFilmSequelRepository(storage *postgres.Storage) *filmSequelRepository {
	return &filmSequelRepository{
		storage: storage,
	}
}

func (s *filmSequelRepository) GetAll(ctx context.Context, filmId string) ([]service.Sequel, error) {
	var filmSequels []service.FilmsSequel
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmId).
		Preload("Sequel").
		Find(&filmSequels)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []service.Sequel{}, nil
	} else if result.Error != nil {
		return []service.Sequel{}, result.Error
	}
	fmt.Println()
	sequels := make([]service.Sequel, len(filmSequels))
	for i, filmSequel := range filmSequels {
		sequels[i] = service.Sequel{
			SequelId:     filmSequel.Sequel.SequelId,
			NameRU:       filmSequel.Sequel.NameRU,
			NameOriginal: filmSequel.Sequel.NameOriginal,
			PosterURL:    filmSequel.Sequel.PosterURL,
		}
	}

	return sequels, nil
}

func (s *filmSequelRepository) Save(sequel []service.Sequel) error {
	return nil
}
