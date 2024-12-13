package userfilm

import (
	"kbox-api/internal/model"
)

func MapDomainUserFilmToResponseDTO(userFilm model.UserFilm) GetUserFilmResponseDTO {
	return GetUserFilmResponseDTO{
		UserID: userFilm.UserID,
		FilmID: userFilm.FilmID,
		Film:   userFilm.Film,
	}
}
