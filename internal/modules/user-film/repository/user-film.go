package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user-film/service"
	"strconv"
)

type userFilmRepository struct {
	storage *postgres.Storage
}

var ErrAlreadyAdded = errors.New("film is already added")
var ErrNotFound = errors.New("film not found")

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
		return []service.UserFilm{}, nil
	} else if result.Error != nil {
		return []service.UserFilm{}, result.Error
	}
	return userFilms, nil
}

func (u *userFilmRepository) Add(ctx context.Context, userId string, filmId string) error {
	filmIdInt, err := strconv.Atoi(filmId)
	if err != nil {
		return err
	}

	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmIdInt).Find(&service.UserFilm{})
	if isFavourite.RowsAffected > 0 {
		return ErrAlreadyAdded
	}
	result := u.storage.DB().WithContext(ctx).Create(&service.UserFilm{
		UserId: userId,
		FilmId: filmIdInt,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userFilmRepository) Delete(ctx context.Context, userId string, filmId string) error {
	isFavourite := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Find(&service.UserFilm{})
	if isFavourite.RowsAffected == 0 {
		return ErrNotFound
	}
	result := u.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Delete(&service.UserFilm{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
