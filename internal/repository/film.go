package repository

import (
	"context"
	"kinopoisk-api/internal/service"
)

type filmRepository struct {
}

func (f *filmRepository) GetOne(ctx context.Context, filmId string) (service.Film, error) {
	panic("implement me")
}

func NewFilmRepository() *filmRepository {
	return &filmRepository{}
}
