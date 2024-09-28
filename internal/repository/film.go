package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
)

type filmRepository struct {
	storage *postgres.Storage
}

func (f *filmRepository) GetOne(ctx context.Context, filmId string) (service.Film, error) {
	var film service.Film

	result := f.storage.DB().WithContext(ctx).Where("film_id = ?", filmId).First(&film)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return service.Film{}, nil
	} else if result.Error != nil {
		return service.Film{}, result.Error
	}

	return film, nil
}

func (f *filmRepository) Save(film service.Film) error {
	result := f.storage.DB().Create(&film)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewFilmRepository(storage *postgres.Storage) *filmRepository {
	return &filmRepository{
		storage: storage,
	}
}
