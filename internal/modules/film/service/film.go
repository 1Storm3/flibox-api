package service

import (
	"context"

	"github.com/lib/pq"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/external/service"
	"kbox-api/internal/modules/film/dto"
	"kbox-api/internal/modules/film/mapper"
	"kbox-api/internal/modules/film/repository"
)

type FilmServiceInterface interface {
	GetOne(filmId string) (dto.FilmResponseDTO, error)
	Search(match string, genres []string, page, pageSize int) ([]dto.FilmSearchResponseDTO, int64, error)
}

type FilmService struct {
	filmRepo        repository.FilmRepositoryInterface
	externalService service.ExternalServiceInterface
}

func NewFilmService(
	filmRepo repository.FilmRepositoryInterface,
	externalService service.ExternalServiceInterface,
) *FilmService {
	return &FilmService{
		filmRepo:        filmRepo,
		externalService: externalService,
	}
}

func (f *FilmService) GetOne(filmId string) (dto.FilmResponseDTO, error) {
	result, err := f.filmRepo.GetOne(context.Background(), filmId)
	if err != nil {
		return dto.FilmResponseDTO{}, err
	}

	if result.ID == nil {
		externalFilm, err := f.externalService.SetExternalRequest(filmId)
		if err != nil {
			return dto.FilmResponseDTO{}, err
		}
		var genres []string
		for _, genre := range externalFilm.Genres {
			genres = append(genres, genre.Genre)
		}

		film := model.Film{
			ID:              externalFilm.ID,
			NameRU:          externalFilm.NameRU,
			NameOriginal:    externalFilm.NameOriginal,
			Year:            externalFilm.Year,
			PosterURL:       externalFilm.PosterURL,
			RatingKinopoisk: externalFilm.RatingKinopoisk,
			Description:     externalFilm.Description,
			LogoURL:         externalFilm.LogoURL,
			Type:            externalFilm.Type,
			Genres:          pq.StringArray(genres),
		}

		if err := f.filmRepo.Save(film); err != nil {
			return dto.FilmResponseDTO{}, err
		}

		filmDTO := mapper.MapModelFilmToResponseDTO(film)

		return filmDTO, nil
	}

	filmDTO := mapper.MapModelFilmToResponseDTO(result)
	return filmDTO, nil
}

func (f *FilmService) Search(match string, genres []string, page int, pageSize int) ([]dto.FilmSearchResponseDTO, int64, error) {
	films, totalRecords, err := f.filmRepo.Search(match, genres, page, pageSize)

	if err != nil {
		return []dto.FilmSearchResponseDTO{}, 0, err
	}

	var filmsDTO []dto.FilmSearchResponseDTO
	for _, film := range films {
		filmsDTO = append(filmsDTO, mapper.MapModelFilmToResponseSearchDTO(film))
	}

	return filmsDTO, totalRecords, nil
}
