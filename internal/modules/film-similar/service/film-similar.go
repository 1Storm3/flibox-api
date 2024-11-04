package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"kbox-api/internal/config"
	"kbox-api/internal/modules/film-similar/dto"
	dtoFilm "kbox-api/internal/modules/film/dto"
	"kbox-api/internal/modules/film/handler"
	"kbox-api/shared/httperror"
)

type FilmSimilarService struct {
	filmSimilarRepo FilmSimilarRepository
	filmService     handler.FilmService
	config          *config.Config
}

const baseUrlForAllSimilar = "https://kinopoiskapiunofficial.tech/api/v2.2/films/%s/similars"

func NewFilmsSimilarService(
	filmSimilarRepo FilmSimilarRepository,
	config *config.Config,
	filmService handler.FilmService,
) *FilmSimilarService {
	return &FilmSimilarService{
		filmSimilarRepo: filmSimilarRepo,
		filmService:     filmService,
		config:          config,
	}
}

func (s *FilmSimilarService) GetAll(filmId string) ([]dtoFilm.FilmResponseDTO, error) {
	result, err := s.filmSimilarRepo.GetAll(context.Background(), filmId)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, err
	}

	if len(result) > 0 {
		var similars []dtoFilm.FilmResponseDTO
		for _, similar := range result {
			res, err := s.filmService.GetOne(strconv.Itoa(similar.SimilarId))

			if err != nil {
				return []dtoFilm.FilmResponseDTO{}, err
			}
			similars = append(similars, res)
		}
		return similars, nil
	}
	similars, err := s.FetchSimilar(filmId)
	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}
	return similars, nil
}

func (s *FilmSimilarService) FetchSimilar(filmId string) ([]dtoFilm.FilmResponseDTO, error) {
	apikey := s.config.DB.ApiKey
	baseUrlForAllSimilar := fmt.Sprintf(baseUrlForAllSimilar, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSimilar, nil)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	req.Header.Add("X-API-KEY", apikey)

	client := &http.Client{}
	resAllSimilars, err := client.Do(req)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{},
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
		return []dtoFilm.FilmResponseDTO{}, httperror.New(
			http.StatusConflict,
			"Код ответа Kinopoisk API: "+resAllSimilars.Status,
		)
	}
	bodyAllSimilars, err := io.ReadAll(resAllSimilars.Body)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var apiResponse struct {
		Total int                         `json:"total"`
		Items []dto.GetExternalSimilarDTO `json:"items"`
	}

	if apiResponse.Total == 0 || len(apiResponse.Items) == 0 {
		return nil, httperror.New(
			http.StatusNotFound,
			"Похожие фильмы не найдены",
		)
	}

	err = json.Unmarshal(bodyAllSimilars, &apiResponse)
	var similars []dtoFilm.FilmResponseDTO
	for _, similar := range apiResponse.Items {
		film, err := s.filmService.GetOne(strconv.Itoa(similar.FilmId))

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []dtoFilm.FilmResponseDTO{},
				httperror.New(
					http.StatusInternalServerError,
					err.Error(),
				)
		}

		err = s.filmSimilarRepo.Save(filmIdInt, similar.FilmId)
		if err != nil {
			return nil, err
		}

		similars = append(similars, film)
	}
	return similars, nil
}
