package film

import "github.com/lib/pq"

type SearchResponseDTO struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	PosterURL       *string  `json:"posterUrl"`
}

type ResponseDTO struct {
	ID              *int           `json:"kinopoiskId"`
	NameRU          *string        `json:"nameRu"`
	NameOriginal    *string        `json:"nameOriginal"`
	Year            *int           `json:"year"`
	RatingKinopoisk *float64       `json:"ratingKinopoisk"`
	PosterURL       *string        `json:"posterUrl"`
	Description     *string        `json:"description"`
	LogoURL         *string        `json:"logoUrl"`
	Type            *string        `json:"type"`
	CoverURL        *string        `json:"coverUrl"`
	TrailerURL      *string        `json:"trailerUrl"`
	Genres          pq.StringArray `json:"genres"`
}
