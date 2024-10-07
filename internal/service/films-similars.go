package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/config"
	"net/http"
	"strconv"
)

type ExternalSimilar struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}

type FilmSimilar struct {
	SimilarId int  `json:"similarId" gorm:"column:similar_id"`
	FilmId    int  `json:"filmId" gorm:"column:film_id"`
	Film      Film `gorm:"foreignKey:FilmId;references:ID"`
}

type FilmSimilarService struct {
	filmSimilarRepo FilmSimilarRepository
	filmService     FilmServiceI
	config          *config.Config
}

const baseUrlForAllSimilar = "https://kinopoiskapiunofficial.tech/api/v2.2/films/%s/similars"

func NewFilmsSimilarService(filmSimilarRepo FilmSimilarRepository, config *config.Config, filmService FilmServiceI) *FilmSimilarService {
	return &FilmSimilarService{
		filmSimilarRepo: filmSimilarRepo,
		filmService:     filmService,
		config:          config,
	}
}

func (s *FilmSimilarService) GetAll(filmId string) ([]Film, error) {
	result, err := s.filmSimilarRepo.GetAll(context.Background(), filmId)

	if err != nil {
		return []Film{}, fmt.Errorf("failed to fetch similar from repository: %w", err)
	}

	if len(result) > 0 {
		var similars []Film
		for _, similar := range result {
			res, err := s.filmService.GetOne(strconv.Itoa(similar.FilmId))

			if err != nil {
				return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
			}
			similars = append(similars, res)
		}
		return similars, nil
	}
	similars, err := s.FetchSimilar(filmId)
	if err != nil {
		return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}
	return similars, nil
}

func (s *FilmSimilarService) FetchSimilar(filmId string) ([]Film, error) {
	apikey := s.config.DB.ApiKey
	baseUrlForAllSimilar := fmt.Sprintf(baseUrlForAllSimilar, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSimilar, nil)

	if err != nil {
		return []Film{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-API-KEY", apikey)

	client := &http.Client{}
	resAllSimilars, err := client.Do(req)

	if err != nil {
		return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}

	defer resAllSimilars.Body.Close()

	if resAllSimilars.StatusCode != 200 {
		return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}
	bodyAllSimilars, err := io.ReadAll(resAllSimilars.Body)

	var apiResponse struct {
		Total int               `json:"total"`
		Items []ExternalSimilar `json:"items"`
	}

	err = json.Unmarshal(bodyAllSimilars, &apiResponse)
	if err != nil {
		return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
	}
	var similars []Film
	for _, similar := range apiResponse.Items {
		film, err := s.filmService.GetOne(strconv.Itoa(similar.FilmId))

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []Film{}, fmt.Errorf("failed to fetch similar from Kinopoisk API: %w", err)
		}

		err = s.filmSimilarRepo.Save(filmIdInt, similar.FilmId)

		if err != nil {
			return nil, err
		}

		similars = append(similars, film)
	}
	return similars, nil
}
