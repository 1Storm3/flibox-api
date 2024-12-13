package collection

import (
	"context"
	"errors"
	"kbox-api/internal/shared/httperror"
	"net/http"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
)

type RepositoryInterface interface {
	GetOne(ctx context.Context, collectionId string) (model.Collection, error)
	Create(ctx context.Context, collection model.Collection) (model.Collection, error)
	GetAll(ctx context.Context, page, pageSize int) ([]model.Collection, int64, error)
	Delete(ctx context.Context, collectionId string) error
	Update(ctx context.Context, collection model.Collection, collectionId string) (model.Collection, error)
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]model.Collection, int64, error)
}

type Repository struct {
	storage *postgres.Storage
}

func NewCollectionRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (c *Repository) GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]model.Collection, int64, error) {
	var collections []model.Collection
	var totalRecords int64
	err := c.storage.DB().WithContext(ctx).Model(&model.Collection{}).Where("user_id = ?", userID).Count(&totalRecords).Error
	if err != nil {
		return []model.Collection{}, 0, err
	}
	err = c.storage.DB().WithContext(ctx).Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&collections).Error
	if err != nil {
		return []model.Collection{}, 0, err
	}
	return collections, totalRecords, nil
}

func (c *Repository) GetAll(ctx context.Context, page, pageSize int) ([]model.Collection, int64, error) {
	var collections []model.Collection
	var totalRecords int64
	err := c.storage.DB().WithContext(ctx).Model(&model.Collection{}).Count(&totalRecords).Error
	if err != nil {
		return []model.Collection{}, 0, err
	}
	err = c.storage.DB().WithContext(ctx).Preload("User").Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&collections).Error
	if err != nil {
		return []model.Collection{}, 0, err
	}
	return collections, totalRecords, nil
}

func (c *Repository) Delete(ctx context.Context, collectionId string) error {
	err := c.storage.DB().WithContext(ctx).
		Where("id = ?", collectionId).
		Delete(&model.Collection{}).Error

	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return nil
}

func (c *Repository) Update(ctx context.Context, collection model.Collection, collectionId string) (model.Collection, error) {
	err := c.storage.DB().WithContext(ctx).Model(&collection).Where("id = ?", collectionId).Updates(collection).Error
	if err != nil {
		return model.Collection{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	err = c.storage.DB().WithContext(ctx).Preload("User").First(&collection, "id = ?", collectionId).Error
	if err != nil {
		return model.Collection{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return collection, nil
}

func (c *Repository) GetOne(ctx context.Context, collectionId string) (model.Collection, error) {
	var collection model.Collection
	err := c.storage.DB().WithContext(ctx).
		Where("id = ?", collectionId).
		Preload("User").
		First(&collection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Collection{}, httperror.New(
				http.StatusNotFound,
				"Коллекция не найдена",
			)
		}
		return model.Collection{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return collection, nil
}

func (c *Repository) Create(ctx context.Context, collection model.Collection) (model.Collection, error) {
	result := c.storage.DB().WithContext(ctx).Create(&collection)

	if result.Error != nil {
		return model.Collection{}, result.Error
	}
	return collection, nil
}
