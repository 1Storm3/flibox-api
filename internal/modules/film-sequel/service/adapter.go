package service

import "context"

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]FilmSequel, error)
	Save(filmId int, sequelId int) error
}
