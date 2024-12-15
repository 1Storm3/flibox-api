package userfilm

import "github.com/1Storm3/flibox-api/internal/model"

type GetUserFilmResponseDTO struct {
	UserID string             `json:"userId"`
	FilmID int                `json:"filmId"`
	Type   model.TypeUserFilm `json:"type"`
	Film   model.Film
}

type Params struct {
	UserID string
	FilmID string
	Type   model.TypeUserFilm
}
