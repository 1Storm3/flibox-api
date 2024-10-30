package handler

import filmservice "kinopoisk-api/internal/modules/film/service"

type FilmService interface {
	GetOne(filmId string) (filmservice.Film, error)
	Search(match string, genres []string, page, pageSize int) ([]filmservice.FilmSearch, int64, error)
}
