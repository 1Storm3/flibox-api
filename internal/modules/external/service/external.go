package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kinopoisk-api/internal/config"
	"kinopoisk-api/shared/httperror"
)

const baseUrlForAllFilms = "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

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

type Genre struct {
	Genre string `json:"genre"`
}

type ExternalService struct {
	config *config.Config
}

func NewExternalService(config *config.Config) *ExternalService {
	return &ExternalService{
		config: config,
	}
}

func (e *ExternalService) SetExternalRequest(filmId string) (ExternalFilm, error) {
	apiKey := e.config.DB.ApiKey
	urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)
	req, err := http.NewRequest("GET", urlAllFilms, nil)
	if err != nil {
		return ExternalFilm{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllFilms, err := client.Do(req)
	if err != nil {
		return ExternalFilm{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resAllFilms.Body)

	if resAllFilms.StatusCode == http.StatusNotFound {
		return ExternalFilm{}, httperror.New(http.StatusNotFound, "Фильм не найден")
	}

	if resAllFilms.StatusCode != http.StatusOK {
		return ExternalFilm{},
			httperror.New(
				http.StatusInternalServerError,
				"Не удалось получить данные о фильме c внешнего апи"+resAllFilms.Status,
			)
	}

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return ExternalFilm{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalFilm ExternalFilm
	err = json.Unmarshal(bodyAllFilms, &externalFilm)
	if err != nil {
		return ExternalFilm{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	return externalFilm, nil
}
