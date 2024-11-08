package dto

import "kbox-api/internal/model"

type GetUserFilmResponseDTO struct {
	UserID string `json:"userId"`
	FilmID int    `json:"filmId"`
	Film   model.Film
}
