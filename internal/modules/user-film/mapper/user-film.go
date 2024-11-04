package mapper

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/user-film/dto"
)

func MapDomainUserFilmToResponseDTO(userFilm model.UserFilm) dto.GetUserFilmResponseDTO {
	return dto.GetUserFilmResponseDTO{
		UserId: userFilm.UserId,
		FilmId: userFilm.FilmId,
		Film:   userFilm.Film,
	}
}
