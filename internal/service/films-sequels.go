package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kinopoisk-api/internal/config"
	"net/http"
)

type FilmsSequel struct {
	SequelId int    `json:"sequelId" gorm:"column:sequel_id"`
	FilmId   int    `json:"filmId" gorm:"column:film_id"`
	Sequel   Sequel `gorm:"foreignKey:SequelId;references:SequelId"`
}

type FilmSequelService struct {
	filmSequelRepo FilmSequelRepository
	sequelRepo     SequelRepository
	config         *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmsSequelService(filmSequelRepo FilmSequelRepository, sequelRepo SequelRepository, config *config.Config) *FilmSequelService {
	return &FilmSequelService{
		filmSequelRepo: filmSequelRepo,
		config:         config,
		sequelRepo:     sequelRepo,
	}
}

func (s *FilmSequelService) GetAll(filmId string) ([]Sequel, error) {
	result, err := s.filmSequelRepo.GetAll(context.Background(), filmId)

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

func (s *FilmSequelService) FetchSequels(filmId string) ([]Sequel, error) {
	apiKey := s.config.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

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
		return []Sequel{}, fmt.Errorf("kinopoisk API request failed with status: %d", resAllSequels.StatusCode)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []Sequel{}, fmt.Errorf("failed to read response body: %w", err)
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
		return []Sequel{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return sequels, nil
}
