package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

var _ ServiceInterface = (*Service)(nil)

const baseUrlForAllFilms = "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

type ServiceInterface interface {
	SetExternalRequest(filmId string) (GetExternalFilmDTO, error)
}

type Service struct {
	cfg *config.Config
}

func NeewExternalService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) SetExternalRequest(filmId string) (GetExternalFilmDTO, error) {
	apiKey := s.cfg.DB.ApiKey
	urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)
	req, err := http.NewRequest("GET", urlAllFilms, nil)
	if err != nil {
		return GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllFilms, err := client.Do(req)
	if err != nil {
		return GetExternalFilmDTO{},
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
		return GetExternalFilmDTO{}, httperror.New(http.StatusNotFound, "Фильм не найден")
	}

	if resAllFilms.StatusCode != http.StatusOK {
		return GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				"Не удалось получить данные о фильме c внешнего апи"+resAllFilms.Status,
			)
	}

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalFilm GetExternalFilmDTO
	err = json.Unmarshal(bodyAllFilms, &externalFilm)
	if err != nil {
		return GetExternalFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	return externalFilm, nil
}
