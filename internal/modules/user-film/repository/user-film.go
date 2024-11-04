package repository

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/shared/httperror"
)

type userFilmRepository struct {
	storage *postgres.Storage
}

func NewUserFilmRepository(storage *postgres.Storage) *userFilmRepository {
	return &userFilmRepository{
		storage: storage,
	}
}

func (u *userFilmRepository) GetAll(ctx context.Context, userId string) ([]model.UserFilm, error) {
	var userFilms []model.UserFilm
	result := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ?", userId).
		Preload("Film").
		Find(&userFilms)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.UserFilm{},
			httperror.New(
				http.StatusNotFound,
				"Избранные фильмы не найдены у этого пользователя",
			)
	} else if result.Error != nil {
		return []model.UserFilm{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return userFilms, nil
}

func (u *userFilmRepository) Add(ctx context.Context, userId string, filmId string) error {
	filmIdInt, err := strconv.Atoi(filmId)
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			"Неверный формат ID фильма",
		)
	}

	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmIdInt).Find(&model.UserFilm{})
	if isFavourite.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в избранное",
		)
	}
	result := u.storage.DB().WithContext(ctx).Create(&model.UserFilm{
		UserId: userId,
		FilmId: filmIdInt,
	})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}

func (u *userFilmRepository) Delete(ctx context.Context, userId string, filmId string) error {
	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Find(&model.UserFilm{})
	if isFavourite.RowsAffected == 0 {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден в избранном",
		)
	}
	result := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Delete(&model.UserFilm{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}
