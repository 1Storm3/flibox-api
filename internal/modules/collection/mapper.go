package collection

import (
	"kbox-api/internal/model"
)

func MapModelCollectionToResponseDTO(collection model.Collection) ResponseDTO {
	return ResponseDTO{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		User: User{
			ID:       collection.User.ID,
			NickName: collection.User.NickName,
			Photo:    collection.User.Photo,
		},
		Tags:      collection.Tags,
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}
}

func MapModelFilmToDTO(film model.Film) Film {
	return Film{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		PosterURL:       film.PosterURL,
		RatingKinopoisk: film.RatingKinopoisk,
	}
}

func MapModelFilmsToDTOs(films []model.Film) []Film {
	dtoFilms := make([]Film, len(films))
	for i, film := range films {
		dtoFilms[i] = MapModelFilmToDTO(film)
	}
	return dtoFilms
}
