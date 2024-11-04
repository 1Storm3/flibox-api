package handler

import (
	"kbox-api/internal/modules/film/dto"
)

type FilmSequelService interface {
	GetAll(filmId string) ([]dto.FilmResponseDTO, error)
}
