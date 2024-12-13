package collectionfilm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/internal/shared/httperror"
)

type RepositoryInterface interface {
	Add(ctx context.Context, collectionId string, filmId int) error
	Delete(ctx context.Context, collectionId string, filmId int) error
	GetFilmsByCollectionId(ctx context.Context, collectionID string, page int, pageSize int) ([]model.Film, int64, error)
}

type Repository struct {
	storage *postgres.Storage
}

func NewCollectionFilmRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (c *Repository) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int, pageSize int,
) ([]model.Film, int64, error) {
	var films []model.Film
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).
		Model(&model.CollectionFilm{}).
		Where("collection_id = ?", collectionID).
		Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = c.storage.DB().WithContext(ctx).
		Model(&model.CollectionFilm{}).
		Select("films.id, films.name_ru, films.name_original, films.year, films.poster_url, films.rating_kinopoisk").
		Joins("JOIN films ON films.id = collection_films.film_id").
		Where("collection_films.collection_id = ?", collectionID).
		Offset(offset).
		Limit(pageSize).
		Scan(&films).Error
	if err != nil {
		return nil, 0, err
	}
	return films, totalRecords, nil
}

func (c *Repository) Add(
	ctx context.Context,
	collectionId string,
	filmId int,
) error {
	var collection model.Collection
	err := c.storage.DB().WithContext(ctx).
		Table("collections").
		Where("id = ?", collectionId).
		First(&collection).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.New(
				http.StatusNotFound,
				"Коллекция не найдена",
			)
		}
		return httperror.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Ошибка при получении коллекции: %v", err),
		)
	}

	newCollectionFilm := &model.CollectionFilm{
		CollectionID: collectionId,
		FilmID:       filmId,
	}
	err = c.storage.DB().WithContext(ctx).Create(newCollectionFilm).Error

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return httperror.New(
				http.StatusConflict,
				"Фильм уже добавлен в коллекцию",
			)
		}
		if strings.Contains(err.Error(), "collection_films_film_id_fkey") {
			return httperror.New(
				http.StatusNotFound,
				"Фильм не найден",
			)
		}
		return httperror.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Ошибка при добавлении фильма в коллекцию: %v", err),
		)
	}

	return nil
}

func (c *Repository) Delete(ctx context.Context, collectionId string, filmId int) error {
	err := c.storage.DB().WithContext(ctx).
		Where("collection_id = ? AND film_id = ?", collectionId, filmId).
		Delete(&model.CollectionFilm{}).Error

	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return nil
}
