package service

import "context"

type UserFilmRepository interface {
	GetAll(ctx context.Context, userId string) ([]UserFilm, error)
	Add(ctx context.Context, userId string, filmId string) error
	Delete(ctx context.Context, userId string, filmId string) error
}
