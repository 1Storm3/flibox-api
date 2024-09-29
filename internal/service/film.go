package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/shared/httperror"
	"net/http"
)

type Film struct {
	ID              *int     `json:"kinopoiskId" gorm:"column:id"`
	NameRU          *string  `json:"nameRu" gorm:"column:name_ru"`
	NameOriginal    *string  `json:"nameOriginal" gorm:"column:name_original"`
	Year            *int     `json:"year" gorm:"column:year"`
	PosterURL       *string  `json:"posterUrl" gorm:"column:poster_url"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	Description     *string  `json:"description" gorm:"column:description"`
	LogoURL         *string  `json:"logoUrl" gorm:"column:logo_url"`
	Type            *string  `json:"type" gorm:"column:type"`
	Sequels         []*Film  `gorm:"many2many:film_sequels;joinForeignKey:film_id;JoinReferences:sequel_id"`
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
		baseUrlForAllFilms := "https://kinopoiskapiunofficial.tech/api/v2.2/films/"
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
		var film Film
		err = json.Unmarshal(bodyAllFilms, &film)

		if err != nil {
			return Film{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		if err := f.filmRepo.Save(film); err != nil {
			return Film{}, fmt.Errorf("failed to save film to repository: %w", err)
		}

		return film, nil
	}

	return result, nil
}
