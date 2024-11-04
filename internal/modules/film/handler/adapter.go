package handler

import (
	"kbox-api/internal/modules/film/dto"
)

type FilmService interface {
	GetOne(filmId string) (dto.FilmResponseDTO, error)
	Search(match string, genres []string, page, pageSize int) ([]dto.FilmSearchResponseDTO, int64, error)
}
