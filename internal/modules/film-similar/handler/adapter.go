package handler

import filmservice "kinopoisk-api/internal/modules/film/service"

type FilmSimilarService interface {
	GetAll(filmId string) ([]filmservice.Film, error)
}
