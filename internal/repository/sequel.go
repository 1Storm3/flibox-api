package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
)

type sequelRepository struct {
	storage *postgres.Storage
}

func (s *sequelRepository) GetAll(ctx context.Context, filmId string) ([]service.Sequel, error) {
	var sequels []service.Sequel

	result := s.storage.DB().WithContext(ctx).Where("film_id = ?", filmId).Find(&sequels)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []service.Sequel{}, nil
	} else if result.Error != nil {
		return []service.Sequel{}, result.Error
	}

	return sequels, nil
}

func (s *sequelRepository) Save(sequel []service.Sequel) error {
	tx := s.storage.DB().Begin()

	result := tx.Create(&sequel)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	var filmsSequels []service.FilmsSequel
	for _, v := range sequel {
		filmsSequels = append(filmsSequels, service.FilmsSequel{
			SequelId: v.SequelId,
			FilmId:   v.FilmId,
		})
	}

	result = tx.Create(&filmsSequels)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func NewSequelRepository(storage *postgres.Storage) *sequelRepository {
	return &sequelRepository{
		storage: storage,
	}
}
