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

func (s *sequelRepository) Save(filmId string, sequels []service.Sequel) error {
	var existingSequels []service.Sequel
	sequelIds := make([]int, len(sequels))
	for i, seq := range sequels {
		sequelIds[i] = seq.SequelId
	}

	err := s.storage.DB().WithContext(context.Background()).
		Where("sequel_id IN ?", sequelIds).
		Find(&existingSequels).Error
	if err != nil {
		return err
	}

	existingMap := make(map[int]bool)
	for _, existing := range existingSequels {
		existingMap[existing.SequelId] = true
	}

	var newSequels []service.Sequel
	for _, sequel := range sequels {
		if !existingMap[sequel.SequelId] {
			newSequels = append(newSequels, sequel)
		}
	}

	if len(newSequels) > 0 {
		tx := s.storage.DB().Begin()

		result := tx.Create(&newSequels)
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
		for _, v := range newSequels {
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

	return nil
}

func NewSequelRepository(storage *postgres.Storage) *sequelRepository {
	return &sequelRepository{
		storage: storage,
	}
}
