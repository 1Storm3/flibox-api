package repository

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/film/service"
)

type filmRepository struct {
	storage *postgres.Storage
}

func (f *filmRepository) GetOne(ctx context.Context, filmId string) (service.Film, error) {
	var film service.Film

	result := f.storage.DB().WithContext(ctx).Where("id = ?", filmId).First(&film)
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

func (f *filmRepository) Search(
	match string,
	genres []string,
	limit, pageSize int,
) ([]service.FilmSearch, int64, error) {
	var films []service.FilmSearch
	var totalRecords int64

	offset := (limit - 1) * pageSize

	query := f.storage.DB().Table("films")

	query = query.Where("name_ru ILIKE ? OR name_original ILIKE ?", "%"+match+"%", "%"+match+"%")

	if len(genres) > 0 {
		query = query.Where("genres && ?", pq.Array(genres))
	}

	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Limit(pageSize).
		Offset(offset).
		Find(&films).Error

	if err != nil {
		return nil, 0, err
	}

	return films, totalRecords, nil
}

func NewFilmRepository(storage *postgres.Storage) *filmRepository {
	return &filmRepository{
		storage: storage,
	}
}
