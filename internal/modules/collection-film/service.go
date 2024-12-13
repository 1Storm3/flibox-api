package collectionfilm

import (
	"context"

	"kbox-api/internal/modules/collection"
)

type ServiceInterface interface {
	Add(ctx context.Context, collectionId string, filmDto CreateCollectionFilmDTO) error
	Delete(ctx context.Context, collectionId string, filmDto DeleteCollectionFilmDTO) error
	GetFilmsByCollectionId(ctx context.Context, collectionID string, page int, pageSize int) (films FilmsByCollectionIdDTO, totalRecords int64, err error)
}

type Service struct {
	repository RepositoryInterface
}

func NewCollectionFilmService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (c *Service) Add(
	ctx context.Context,
	collectionId string,
	filmDto CreateCollectionFilmDTO,
) error {
	return c.repository.Add(ctx, collectionId, filmDto.FilmID)
}

func (c *Service) Delete(
	ctx context.Context,
	collectionId string,
	filmDto DeleteCollectionFilmDTO,
) error {
	return c.repository.Delete(ctx, collectionId, filmDto.FilmID)
}

func (c *Service) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int,
	pageSize int,
) (films FilmsByCollectionIdDTO, totalRecords int64, err error) {
	result, totalRecords, err := c.repository.GetFilmsByCollectionId(ctx, collectionID, page, pageSize)

	return FilmsByCollectionIdDTO{
		CollectionID: collectionID,
		Films:        collection.MapModelFilmsToDTOs(result),
	}, totalRecords, err
}
