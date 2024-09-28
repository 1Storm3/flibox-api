package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
	"strconv"
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

func (s *sequelRepository) Save(filmId string, sequel []service.Sequel) error {
	isExist := s.storage.DB().WithContext(context.Background()).Where("sequel_id = ?", filmId).Find(&sequel)

	if isExist != nil {
		return nil
	}
	tx := s.storage.DB().Begin()

	result := tx.Create(&sequel)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	filmIdInt, err := strconv.Atoi(filmId)

	if err != nil {
		tx.Rollback()
		return err
	}

	var filmsSequels []service.FilmsSequel
	for _, v := range sequel {
		filmsSequels = append(filmsSequels, service.FilmsSequel{
			SequelId: v.SequelId,
			FilmId:   filmIdInt,
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