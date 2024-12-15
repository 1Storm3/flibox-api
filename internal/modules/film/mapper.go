package film

import "github.com/1Storm3/flibox-api/internal/model"

func MapModelFilmToResponseSearchDTO(film model.Film) SearchResponseDTO {
	return SearchResponseDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
	}
}

func MapModelFilmToResponseDTO(film model.Film) ResponseDTO {
	return ResponseDTO{
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
		CoverURL:        film.CoverURL,
		TrailerURL:      film.TrailerURL,
	}
}
