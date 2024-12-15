package film

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/modules/external"
	"github.com/lib/pq"
)

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	GetOne(ctx context.Context, filmId string) (ResponseDTO, error)
	Search(ctx context.Context, match string, genres []string, page, pageSize int) ([]SearchResponseDTO, int64, error)
	GetOneByNameRu(ctx context.Context, nameRu string) (ResponseDTO, error)
}

type Service struct {
	repository      RepositoryInterface
	externalService external.ServiceInterface
}

func NewFilmService(
	repository RepositoryInterface,
	externalService external.ServiceInterface,
) *Service {
	return &Service{
		repository:      repository,
		externalService: externalService,
	}
}

func (f *Service) GetOne(ctx context.Context, filmId string) (ResponseDTO, error) {
	result, err := f.repository.GetOne(ctx, filmId)
	if err != nil {
		return ResponseDTO{}, err
	}

	if result.ID == nil {
		externalFilm, err := f.externalService.SetExternalRequest(filmId)
		if err != nil {
			return ResponseDTO{}, err
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

		if err := f.repository.Save(ctx, film); err != nil {
			return ResponseDTO{}, err
		}

		filmDTO := MapModelFilmToResponseDTO(film)

		return filmDTO, nil
	}

	filmDTO := MapModelFilmToResponseDTO(result)
	return filmDTO, nil
}

func (f *Service) Search(ctx context.Context, match string, genres []string, page int, pageSize int) ([]SearchResponseDTO, int64, error) {
	films, totalRecords, err := f.repository.Search(ctx, match, genres, page, pageSize)

	if err != nil {
		return []SearchResponseDTO{}, 0, err
	}

	var filmsDTO []SearchResponseDTO
	for _, film := range films {
		filmsDTO = append(filmsDTO, MapModelFilmToResponseSearchDTO(film))
	}

	return filmsDTO, totalRecords, nil
}

func (f *Service) GetOneByNameRu(ctx context.Context, nameRu string) (ResponseDTO, error) {
	result, err := f.repository.GetOneByNameRu(ctx, nameRu)
	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelFilmToResponseDTO(result), nil
}
