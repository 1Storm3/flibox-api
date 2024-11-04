package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kbox-api/internal/config"
	"kbox-api/internal/modules/external/dto"
	"kbox-api/shared/httperror"
)

const baseUrlForAllFilms = "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

type ExternalService struct {
	config *config.Config
}

func NewExternalService(config *config.Config) *ExternalService {
	return &ExternalService{
		config: config,
	}
}

func (e *ExternalService) SetExternalRequest(filmId string) (dto.GetExternalFilmDTO, error) {
	apiKey := e.config.DB.ApiKey
	urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)
	req, err := http.NewRequest("GET", urlAllFilms, nil)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllFilms, err := client.Do(req)
	if err != nil {
		return dto.GetExternalFilmDTO{},
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
		return dto.GetExternalFilmDTO{}, httperror.New(http.StatusNotFound, "Фильм не найден")
	}

	if resAllFilms.StatusCode != http.StatusOK {
		return dto.GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				"Не удалось получить данные о фильме c внешнего апи"+resAllFilms.Status,
			)
	}

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalFilm dto.GetExternalFilmDTO
	err = json.Unmarshal(bodyAllFilms, &externalFilm)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	return externalFilm, nil
}
