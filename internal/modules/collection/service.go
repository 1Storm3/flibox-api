package collection

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/model"
)

type ServiceInterface interface {
	Create(ctx context.Context, collection CreateCollectionDTO, userID string) (ResponseDTO, error)
	GetAll(ctx context.Context, page, pageSize int) ([]ResponseDTO, int64, error)
	GetOne(ctx context.Context, collectionId string) (ResponseDTO, error)
	Update(ctx context.Context, collection UpdateCollectionDTO, collectionId string) (ResponseDTO, error)
	Delete(ctx context.Context, collectionId string) error
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]ResponseDTO, int64, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewCollectionService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (c *Service) Update(ctx context.Context, collection UpdateCollectionDTO, collectionId string) (ResponseDTO, error) {
	result, err := c.repository.Update(ctx, model.Collection{
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		Tags:        collection.Tags,
	}, collectionId)
	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelCollectionToResponseDTO(result), nil
}

func (c *Service) Delete(ctx context.Context, collectionId string) error {
	return c.repository.Delete(ctx, collectionId)
}

func (c *Service) Create(ctx context.Context, collection CreateCollectionDTO, userID string) (ResponseDTO, error) {
	result, err := c.repository.Create(ctx, model.Collection{
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		Tags:        collection.Tags,
		UserId:      &userID,
	})
	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelCollectionToResponseDTO(result), nil
}

func (c *Service) GetAll(ctx context.Context, page, pageSize int) ([]ResponseDTO, int64, error) {
	result, totalRecords, err := c.repository.GetAll(ctx, page, pageSize)
	if err != nil {
		return []ResponseDTO{}, 0, err
	}
	var resultDTO []ResponseDTO
	for _, collection := range result {
		resultDTO = append(resultDTO, MapModelCollectionToResponseDTO(collection))
	}
	return resultDTO, totalRecords, nil
}

func (c *Service) GetOne(ctx context.Context, collectionId string) (ResponseDTO, error) {
	result, err := c.repository.GetOne(ctx, collectionId)
	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelCollectionToResponseDTO(result), nil
}

func (c *Service) GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]ResponseDTO, int64, error) {
	result, totalRecords, err := c.repository.GetAllMy(ctx, page, pageSize, userID)
	if err != nil {
		return []ResponseDTO{}, 0, err
	}
	var resultDTO []ResponseDTO
	for _, collection := range result {
		resultDTO = append(resultDTO, MapModelCollectionToResponseDTO(collection))
	}
	return resultDTO, totalRecords, nil
}
