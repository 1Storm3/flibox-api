package handler

import userfilmservice "kinopoisk-api/internal/modules/user-film/service"

type UserFilmService interface {
	GetAll(userId string) ([]userfilmservice.UserFilm, error)
	Add(userId string, filmId string) error
	Delete(userId string, filmId string) error
}
