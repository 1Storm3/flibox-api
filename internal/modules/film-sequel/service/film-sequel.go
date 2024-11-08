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
	"kbox-api/internal/modules/film-sequel/repository"
	dtoFilm "kbox-api/internal/modules/film/dto"
	"kbox-api/internal/modules/film/service"
	"kbox-api/shared/httperror"
)

type FilmSequelServiceInterface interface {
	GetAll(filmId string) ([]dtoFilm.FilmResponseDTO, error)
}

type FilmSequelService struct {
	filmSequelRepo repository.FilmSequelRepositoryInterface
	filmService    service.FilmServiceInterface
	cfg            *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmsSequelService(
	filmSequelRepo repository.FilmSequelRepositoryInterface,
	cfg *config.Config,
	filmService service.FilmServiceInterface,
) FilmSequelServiceInterface {
	return &FilmSequelService{
		filmSequelRepo: filmSequelRepo,
		filmService:    filmService,
		cfg:            cfg,
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
			res, err := s.filmService.GetOne(strconv.Itoa(sequel.SequelID))

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
	apiKey := s.cfg.DB.ApiKey
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
