package service

import (
	externalservice "kinopoisk-api/pkg/external-service"
)

type ExternalService interface {
	SetExternalRequest(filmId string) (externalservice.ExternalFilm, error)
}
