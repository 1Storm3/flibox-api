package mapper

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/film/dto"
)

func MapModelFilmToResponseSearchDTO(film model.Film) dto.FilmSearchResponseDTO {
	return dto.FilmSearchResponseDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
	}
}

func MapModelFilmToResponseDTO(film model.Film) dto.FilmResponseDTO {
	return dto.FilmResponseDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
		Description:     film.Description,
		LogoURL:         film.LogoURL,
		Type:            film.Type,
		Genres:          film.Genres,
	}
}
