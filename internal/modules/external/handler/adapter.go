package handler

import (
	externalservice "kinopoisk-api/internal/modules/external/service"
)

type ExternalService interface {
	SetExternalRequest(filmId string) (externalservice.ExternalFilm, error)
}
