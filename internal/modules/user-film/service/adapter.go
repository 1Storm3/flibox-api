package service

import (
	"context"

	"kbox-api/internal/model"
)

type UserFilmRepository interface {
	GetAll(ctx context.Context, userId string) ([]model.UserFilm, error)
	Add(ctx context.Context, userId string, filmId string) error
	Delete(ctx context.Context, userId string, filmId string) error
}
