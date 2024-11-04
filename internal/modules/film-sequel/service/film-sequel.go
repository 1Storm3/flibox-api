package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"kbox-api/internal/config"
	"kbox-api/internal/modules/film-sequel/dto"
	dtoFilm "kbox-api/internal/modules/film/dto"
	"kbox-api/internal/modules/film/handler"
	"kbox-api/shared/httperror"
)

type FilmSequelService struct {
	filmSequelRepo FilmSequelRepository
	filmService    handler.FilmService
	config         *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmsSequelService(
	filmSequelRepo FilmSequelRepository,
	config *config.Config,
	filmService handler.FilmService,
) *FilmSequelService {
	return &FilmSequelService{
		filmSequelRepo: filmSequelRepo,
		filmService:    filmService,
		config:         config,
	}
}

func (s *FilmSequelService) GetAll(filmId string) ([]dtoFilm.FilmResponseDTO, error) {
	result, err := s.filmSequelRepo.GetAll(context.Background(), filmId)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, err
	}
	if len(result) > 0 {
		var sequels []dtoFilm.FilmResponseDTO
		for _, sequel := range result {
			res, err := s.filmService.GetOne(strconv.Itoa(sequel.SequelId))

			if err != nil {
				return []dtoFilm.FilmResponseDTO{}, err
			}
			sequels = append(sequels, res)
		}
		return sequels, nil
	}

	sequels, err := s.FetchSequels(filmId)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, err
	}

	return sequels, nil
}

func (s *FilmSequelService) FetchSequels(filmId string) ([]dtoFilm.FilmResponseDTO, error) {
	apiKey := s.config.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

	if err != nil {
		return []dtoFilm.FilmResponseDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllSequels, err := client.Do(req)
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
	}(resAllSequels.Body)

	if resAllSequels.StatusCode != http.StatusOK {
		return []dtoFilm.FilmResponseDTO{}, httperror.New(
			http.StatusNotFound,
			"Сиквелы не найдены",
		)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []dtoFilm.FilmResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalSequels []dto.GetExternalSequelResponseDTO

	err = json.Unmarshal(bodyAllSequels, &externalSequels)
	var sequels []dtoFilm.FilmResponseDTO
	for _, sequel := range externalSequels {
		film, err := s.filmService.GetOne(strconv.Itoa(sequel.FilmId))

		if err != nil {
			return []dtoFilm.FilmResponseDTO{}, err
		}

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []dtoFilm.FilmResponseDTO{},
				httperror.New(
					http.StatusInternalServerError,
					err.Error(),
				)
		}

		err = s.filmSequelRepo.Save(filmIdInt, sequel.FilmId)
		if err != nil {
			return nil, err
		}

		sequels = append(sequels, film)
	}

	if err != nil {
		return []dtoFilm.FilmResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	return sequels, nil
}
