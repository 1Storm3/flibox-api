package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/modules/film/service"
	"net/http"
	"strconv"

	"kinopoisk-api/internal/config"
)

type ExternalSequel struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}

type FilmSequel struct {
	SequelId int          `json:"sequelId" gorm:"column:sequel_id"`
	FilmId   int          `json:"filmId" gorm:"column:film_id"`
	Film     service.Film `gorm:"foreignKey:FilmId;references:ID"`
}

type FilmSequelService struct {
	filmSequelRepo FilmSequelRepository
	filmService    service.FilmServiceI
	config         *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmsSequelService(
	filmSequelRepo FilmSequelRepository,
	config *config.Config,
	filmService service.FilmServiceI,
) *FilmSequelService {
	return &FilmSequelService{
		filmSequelRepo: filmSequelRepo,
		filmService:    filmService,
		config:         config,
	}
}

func (s *FilmSequelService) GetAll(filmId string) ([]service.Film, error) {
	result, err := s.filmSequelRepo.GetAll(context.Background(), filmId)

	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to fetch sequel from repository: %w", err)
	}
	if len(result) > 0 {
		var sequels []service.Film
		for _, sequel := range result {
			res, err := s.filmService.GetOne(strconv.Itoa(sequel.SequelId))

			if err != nil {
				return []service.Film{}, fmt.Errorf("failed to fetch sequel from Kinopoisk API: %w", err)
			}
			sequels = append(sequels, res)
		}
		return sequels, nil
	}

	sequels, err := s.FetchSequels(filmId)

	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to fetch sequel from Kinopoisk API: %w", err)
	}

	return sequels, nil
}

func (s *FilmSequelService) FetchSequels(filmId string) ([]service.Film, error) {
	apiKey := s.config.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllSequels, err := client.Do(req)
	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to make request to Kinopoisk API: %w", err)
	}
	defer resAllSequels.Body.Close()

	if resAllSequels.StatusCode != http.StatusOK {
		return []service.Film{}, fmt.Errorf("kinopoisk API request failed with status: %d", resAllSequels.StatusCode)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var externalSequels []ExternalSequel
	err = json.Unmarshal(bodyAllSequels, &externalSequels)
	var sequels []service.Film
	for _, sequel := range externalSequels {
		film, err := s.filmService.GetOne(strconv.Itoa(sequel.FilmId))

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []service.Film{}, fmt.Errorf("failed to convert filmId to int: %w", err)
		}

		err = s.filmSequelRepo.Save(filmIdInt, sequel.FilmId)
		if err != nil {
			return nil, err
		}

		sequels = append(sequels, film)
	}

	if err != nil {
		return []service.Film{}, fmt.Errorf("failed to decode response from Kinopoisk API: %w", err)
	}
	return sequels, nil
}
