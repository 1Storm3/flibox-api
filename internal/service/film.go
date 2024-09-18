package service

import (
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/config"
	"net/http"
)

type Film struct {
	ExternalId      int     `json:"kinopoiskId"`
	NameRu          string  `json:"nameRu"`
	NameOriginal    string  `json:"nameOriginal"`
	Year            int     `json:"year"`
	PosterUrl       string  `json:"posterUrl"`
	RatingKinopoisk float64 `json:"ratingKinopoisk"`
	Description     string  `json:"description"`
	LogoUrl         string  `json:"logoUrl"`
	Type            string  `json:"type"`
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

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return Film{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var film Film
	err = json.Unmarshal(bodyAllFilms, &film)
	if err != nil {
		return Film{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return film, nil
}
