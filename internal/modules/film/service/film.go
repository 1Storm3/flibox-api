package service

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"kinopoisk-api/internal/service"
)

type Film struct {
	ID              *int           `json:"kinopoiskId" gorm:"column:id"`
	NameRU          *string        `json:"nameRu" gorm:"column:name_ru"`
	NameOriginal    *string        `json:"nameOriginal" gorm:"column:name_original"`
	Year            *int           `json:"year" gorm:"column:year"`
	PosterURL       *string        `json:"posterUrl" gorm:"column:poster_url"`
	RatingKinopoisk *float64       `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	Description     *string        `json:"description" gorm:"column:description"`
	LogoURL         *string        `json:"logoUrl" gorm:"column:logo_url"`
	Type            *string        `json:"type" gorm:"column:type"`
	Sequels         []*Film        `gorm:"many2many:film_sequels;joinForeignKey:film_id;JoinReferences:sequel_id"`
	Similars        []*Film        `gorm:"many2many:film_similars;joinForeignKey:film_id;JoinReferences:similar_id"`
	Genres          pq.StringArray `json:"genres" gorm:"type:text[];column:genres"`
}

type FilmSearch struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	PosterURL       *string  `json:"posterUrl"`
}

type FilmService struct {
	filmRepo        FilmRepository
	externalService service.ExternalService
}

func NewFilmService(
	filmRepo FilmRepository,
	externalService service.ExternalService,
) *FilmService {
	return &FilmService{
		filmRepo:        filmRepo,
		externalService: externalService,
	}
}

func (f *FilmService) GetOne(filmId string) (Film, error) {
	result, err := f.filmRepo.GetOne(context.Background(), filmId)
	if err != nil {
		return Film{}, fmt.Errorf("failed to fetch film from repository: %w", err)
	}

	if result.ID == nil {
		externalFilm, err := f.externalService.SetExternalRequest(filmId)
		if err != nil {
			return Film{}, fmt.Errorf("failed to get film from Kinopoisk API: %w", err)
		}
		var genres []string
		for _, genre := range externalFilm.Genres {
			genres = append(genres, genre.Genre)
		}

		film := Film{
			ID:              externalFilm.ID,
			NameRU:          externalFilm.NameRU,
			NameOriginal:    externalFilm.NameOriginal,
			Year:            externalFilm.Year,
			PosterURL:       externalFilm.PosterURL,
			RatingKinopoisk: externalFilm.RatingKinopoisk,
			Description:     externalFilm.Description,
			LogoURL:         externalFilm.LogoURL,
			Type:            externalFilm.Type,
			Genres:          pq.StringArray(genres),
		}

		if err := f.filmRepo.Save(film); err != nil {
			return Film{}, fmt.Errorf("failed to save film to repository: %w", err)
		}

		return film, nil
	}

	return result, nil
}

func (f *FilmService) Search(match string, genres []string, page int, pageSize int) ([]FilmSearch, int64, error) {
	films, totalRecords, err := f.filmRepo.Search(match, genres, page, pageSize)

	if err != nil {
		return []FilmSearch{}, 0, fmt.Errorf("failed to fetch films from repository: %w", err)
	}

	return films, totalRecords, nil
}
