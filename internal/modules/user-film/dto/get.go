package dto

import "kbox-api/internal/model"

type GetUserFilmResponseDTO struct {
	UserId string `json:"userId"`
	FilmId int    `json:"filmId"`
	Film   model.Film
}
