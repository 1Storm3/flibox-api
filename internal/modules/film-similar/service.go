package filmsimilar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/shared/httperror"
)

type ServiceInterface interface {
	GetAll(ctx context.Context, filmId string) ([]film.ResponseDTO, error)
	FetchSimilar(ctx context.Context, filmId string) ([]film.ResponseDTO, error)
}

type Service struct {
	repository  RepositoryInterface
	filmService film.ServiceInterface
	cfg         *config.Config
}

const baseUrlForAllSimilar = "https://kinopoiskapiunofficial.tech/api/v2.2/films/%s/similars"

func NewFilmsSimilarService(
	repository RepositoryInterface,
	cfg *config.Config,
	filmService film.ServiceInterface,
) *Service {
	return &Service{
		repository:  repository,
		filmService: filmService,
		cfg:         cfg,
	}
}

func (s *Service) GetAll(ctx context.Context, filmId string) ([]film.ResponseDTO, error) {
	result, err := s.repository.GetAll(ctx, filmId)

	if err != nil {
		return []film.ResponseDTO{}, err
	}

	if len(result) > 0 {
		var similars []film.ResponseDTO
		for _, similar := range result {
			res, err := s.filmService.GetOne(ctx, strconv.Itoa(similar.SimilarID))

			if err != nil {
				return []film.ResponseDTO{}, err
			}
			similars = append(similars, res)
		}
		return similars, nil
	}
	similars, err := s.FetchSimilar(ctx, filmId)
	if err != nil {
		return []film.ResponseDTO{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}
	return similars, nil
}

func (s *Service) FetchSimilar(ctx context.Context, filmId string) ([]film.ResponseDTO, error) {
	apikey := s.cfg.DB.ApiKey
	baseUrlForAllSimilar := fmt.Sprintf(baseUrlForAllSimilar, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSimilar, nil)

	if err != nil {
		return []film.ResponseDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	req.Header.Add("X-API-KEY", apikey)

	client := &http.Client{}
	resAllSimilars, err := client.Do(req)

	if err != nil {
		return []film.ResponseDTO{},
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
	}(resAllSimilars.Body)

	if resAllSimilars.StatusCode != http.StatusOK {
		return []film.ResponseDTO{}, httperror.New(
			http.StatusConflict,
			"Код ответа Kinopoisk API: "+resAllSimilars.Status,
		)
	}
	bodyAllSimilars, err := io.ReadAll(resAllSimilars.Body)

	if err != nil {
		return []film.ResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var apiResponse struct {
		Total int                     `json:"total"`
		Items []GetExternalSimilarDTO `json:"items"`
	}

	if apiResponse.Total == 0 || len(apiResponse.Items) == 0 {
		return nil, httperror.New(
			http.StatusNotFound,
			"Похожие фильмы не найдены",
		)
	}

	err = json.Unmarshal(bodyAllSimilars, &apiResponse)

	if err != nil {
		return []film.ResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	var similars []film.ResponseDTO
	for _, similar := range apiResponse.Items {
		filmIsExist, err := s.filmService.GetOne(ctx, strconv.Itoa(similar.FilmId))

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []film.ResponseDTO{},
				httperror.New(
					http.StatusInternalServerError,
					err.Error(),
				)
		}

		err = s.repository.Save(ctx, filmIdInt, similar.FilmId)
		if err != nil {
			return nil, err
		}

		similars = append(similars, filmIsExist)
	}
	return similars, nil
}
