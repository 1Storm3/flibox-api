package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/shared/httperror"
	"net/http"
)

const baseUrlForAllFilms = "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

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
	Genres          pq.StringArray `json:"genres" gorm:"type:text[];column:genres"`
}

type Genre struct {
	Genre string `json:"genre"`
}

type ExternalFilm struct {
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

type FilmService struct {
	filmRepo FilmRepository
	config   *config.Config
}

func NewFilmService(filmRepo FilmRepository, config *config.Config) *FilmService {
	return &FilmService{
		filmRepo: filmRepo,
		config:   config,
	}
}

func (f *FilmService) GetOne(filmId string) (Film, error) {
	result, err := f.filmRepo.GetOne(context.Background(), filmId)
	if err != nil {
		return Film{}, fmt.Errorf("failed to fetch film from repository: %w", err)
	}

	if result.ID == nil {
		apiKey := f.config.DB.ApiKey
		urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)

		req, err := http.NewRequest("GET", urlAllFilms, nil)
		if err != nil {
			return Film{}, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Add("X-API-KEY", apiKey)

		client := &http.Client{}
		resAllFilms, err := client.Do(req)
		if err != nil {
			return Film{}, fmt.Errorf("failed to make request to Kinopoisk API: %w", err)
		}
		defer resAllFilms.Body.Close()

		if resAllFilms.StatusCode == http.StatusNotFound {
			return Film{}, httperror.New(http.StatusNotFound, "Фильм не найден")
		}

		if resAllFilms.StatusCode != http.StatusOK {
			return Film{}, fmt.Errorf("kinopoisk API request failed with status: %d", resAllFilms.StatusCode)
		}

		bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
		if err != nil {
			return Film{}, fmt.Errorf("failed to read response body: %w", err)
		}

		var externalFilm ExternalFilm
		err = json.Unmarshal(bodyAllFilms, &externalFilm)
		if err != nil {
			return Film{}, fmt.Errorf("failed to unmarshal response body: %w", err)
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
