package service

import "context"

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (Film, error)
	Save(film Film) error
}

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (User, error)
}

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]FilmSequel, error)
	Save(filmId int, sequelId int) error
}

type FilmServiceI interface {
	GetOne(filmId string) (Film, error)
}
