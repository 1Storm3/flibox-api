package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user-film/service"
	"kinopoisk-api/shared/httperror"
	"net/http"
	"strconv"
)

type userFilmRepository struct {
	storage *postgres.Storage
}

func NewUserFilmRepository(storage *postgres.Storage) *userFilmRepository {
	return &userFilmRepository{
		storage: storage,
	}
}

func (u *userFilmRepository) GetAll(ctx context.Context, userId string) ([]service.UserFilm, error) {
	var userFilms []service.UserFilm
	result := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ?", userId).
		Preload("Film").
		Find(&userFilms)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []service.UserFilm{},
			httperror.New(
				http.StatusNotFound,
				"Избранные фильмы не найдены у этого пользователя",
			)
	} else if result.Error != nil {
		return []service.UserFilm{},
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

	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmIdInt).Find(&service.UserFilm{})
	if isFavourite.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в избранное",
		)
	}
	result := u.storage.DB().WithContext(ctx).Create(&service.UserFilm{
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
	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Find(&service.UserFilm{})
	if isFavourite.RowsAffected == 0 {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден в избранном",
		)
	}
	result := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Delete(&service.UserFilm{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}
