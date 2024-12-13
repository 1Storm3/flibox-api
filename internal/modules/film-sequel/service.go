package filmsequel

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

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	GetAll(ctx context.Context, filmId string) ([]film.ResponseDTO, error)
}

type Service struct {
	repository  RepositoryInterface
	filmService film.ServiceInterface
	cfg         *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmsSequelService(
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
		var sequels []film.ResponseDTO
		for _, sequel := range result {
			res, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.SequelID))

			if err != nil {
				return []film.ResponseDTO{}, err
			}
			sequels = append(sequels, res)
		}
		return sequels, nil
	}

	sequels, err := s.FetchSequels(ctx, filmId)

	if err != nil {
		return []film.ResponseDTO{}, err
	}

	return sequels, nil
}

func (s *Service) FetchSequels(ctx context.Context, filmId string) ([]film.ResponseDTO, error) {
	apiKey := s.cfg.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

	if err != nil {
		return []film.ResponseDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllSequels, err := client.Do(req)
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
	}(resAllSequels.Body)

	if resAllSequels.StatusCode != http.StatusOK {
		return []film.ResponseDTO{}, httperror.New(
			http.StatusNotFound,
			"Сиквелы не найдены",
		)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []film.ResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalSequels []GetExternalSequelResponseDTO

	err = json.Unmarshal(bodyAllSequels, &externalSequels)
	var sequels []film.ResponseDTO
	for _, sequel := range externalSequels {
		filmExist, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.FilmId))

		if err != nil {
			return []film.ResponseDTO{}, err
		}

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []film.ResponseDTO{},
				httperror.New(
					http.StatusInternalServerError,
					err.Error(),
				)
		}

		err = s.repository.Save(ctx, filmIdInt, sequel.FilmId)
		if err != nil {
			return nil, err
		}

		sequels = append(sequels, filmExist)
	}

	if err != nil {
		return []film.ResponseDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	return sequels, nil
}
