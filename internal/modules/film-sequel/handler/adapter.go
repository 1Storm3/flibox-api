package handler

import filmservice "kinopoisk-api/internal/modules/film/service"

type FilmSequelService interface {
	GetAll(filmId string) ([]filmservice.Film, error)
}
