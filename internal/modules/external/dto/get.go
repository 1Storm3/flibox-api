package dto

type GetExternalFilmDTO struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	PosterURL       *string  `json:"posterUrl"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk"`
	Description     *string  `json:"description"`
	LogoURL         *string  `json:"logoUrl"`
	Type            *string  `json:"type"`
	Genres          []Genre  `json:"genres"`
}

type Genre struct {
	Genre string `json:"genre"`
}
