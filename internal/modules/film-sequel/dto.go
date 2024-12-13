package filmsequel

type GetExternalSequelResponseDTO struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}
