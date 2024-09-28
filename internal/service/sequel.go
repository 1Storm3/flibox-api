package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/config"
	"net/http"
)

type Sequel struct {
	SequelId     int     `json:"sequelId" gorm:"column:sequel_id"`
	NameRU       *string `json:"nameRu" gorm:"column:name_ru"`
	NameOriginal *string `json:"name" gorm:"column:name_original"`
	PosterURL    *string `json:"posterUrl" gorm:"column:poster_url"`
}

type ExternalSequel struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}

type SequelService struct {
	sequelRepo SequelRepository
	config     *config.Config
}

func NewSequelService(sequelRepo SequelRepository, config *config.Config) *SequelService {
	return &SequelService{
		sequelRepo: sequelRepo,
		config:     config,
	}
}

func (s *SequelService) GetAll(filmId string) ([]Sequel, error) {
	result, err := s.sequelRepo.GetAll(context.Background(), filmId)

	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to fetch sequel from repository: %w", err)
	}

	if len(result) > 0 {
		return result, nil
	}

	sequels, err := s.FetchSequels(filmId)

	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to fetch sequel from Kinopoisk API: %w", err)
	}

	return sequels, nil
}

func (s *SequelService) FetchSequels(filmId string) ([]Sequel, error) {
	apiKey := s.config.DB.ApiKey

	urlAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)

	req, err := http.NewRequest("GET", urlAllSequels, nil)

	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}

	resAllSequels, err := client.Do(req)
	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to make request to Kinopoisk API: %w", err)
	}
	defer resAllSequels.Body.Close()

	if resAllSequels.StatusCode != http.StatusOK {
		return []Sequel{}, fmt.Errorf("request to Kinopoisk API failed with status code: %d", resAllSequels.StatusCode)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to read response from Kinopoisk API: %w", err)
	}

	var externalSequels []ExternalSequel
	err = json.Unmarshal(bodyAllSequels, &externalSequels)

	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to decode response from Kinopoisk API: %w", err)
	}
	var sequels []Sequel
	for _, externalSequel := range externalSequels {
		sequel := Sequel{
			SequelId:     externalSequel.FilmId,
			NameRU:       externalSequel.NameRu,
			NameOriginal: externalSequel.NameOriginal,
			PosterURL:    externalSequel.PosterUrl,
		}
		sequels = append(sequels, sequel)
	}

	err = s.sequelRepo.Save(filmId, sequels)

	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to save sequel to repository: %w", err)
	}

	return sequels, nil
}
