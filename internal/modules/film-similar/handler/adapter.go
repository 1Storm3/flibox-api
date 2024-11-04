package handler

import (
	"kbox-api/internal/modules/film/dto"
)

type FilmSimilarService interface {
	GetAll(filmId string) ([]dto.FilmResponseDTO, error)
	FetchSimilar(filmId string) ([]dto.FilmResponseDTO, error)
}
