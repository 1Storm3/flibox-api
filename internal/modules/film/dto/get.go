package dto

import "github.com/lib/pq"

type FilmSearchResponseDTO struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	PosterURL       *string  `json:"posterUrl"`
}

type FilmResponseDTO struct {
	ID              *int           `json:"kinopoiskId"`
	NameRU          *string        `json:"nameRu"`
	NameOriginal    *string        `json:"nameOriginal"`
	Year            *int           `json:"year"`
	RatingKinopoisk *float64       `json:"ratingKinopoisk"`
	PosterURL       *string        `json:"posterUrl"`
	Description     *string        `json:"description"`
	LogoURL         *string        `json:"logoUrl"`
	Type            *string        `json:"type"`
	Genres          pq.StringArray `json:"genres"`
}
