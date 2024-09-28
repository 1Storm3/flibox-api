package service

import "context"

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (Film, error)
	Save(film Film) error
}

type SequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]Sequel, error)
	Save(filmId string, sequel []Sequel) error
}

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (User, error)
}

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]Sequel, error)
	Save(sequel []Sequel) error
}
