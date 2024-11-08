package repository

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/internal/modules/comment/dto"
	"kbox-api/shared/httperror"
)

var _ CommentRepositoryInterface = (*commentRepository)(nil)

type CommentRepositoryInterface interface {
	Create(ctx context.Context, comment model.Comment) (model.Comment, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]model.Comment, int64, error)
	Update(ctx context.Context, comment dto.UpdateCommentDTO, commentID string) (model.Comment, error)
	GetOne(ctx context.Context, commentID string) (model.Comment, error)
}

type commentRepository struct {
	storage *postgres.Storage
}

func NewCommentRepository(storage *postgres.Storage) CommentRepositoryInterface {
	return &commentRepository{
		storage: storage,
	}
}

func (c *commentRepository) Create(ctx context.Context, comment model.Comment) (model.Comment, error) {
	result := c.storage.DB().WithContext(ctx).Create(&comment)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "violates foreign key") {
			return model.Comment{}, httperror.New(
				http.StatusConflict,
				"Родительского комментария не существует с таким ID",
			)
		}
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	return comment, nil
}

func (c *commentRepository) GetOne(ctx context.Context, commentID string) (model.Comment, error) {
	var comment model.Comment

	err := c.storage.DB().WithContext(ctx).Where("id = ?", commentID).First(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Comment{}, httperror.New(
			http.StatusNotFound,
			"Комментарий не найден",
		)
	} else if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comment, nil
}

func (c *commentRepository) Delete(ctx context.Context, commentID string) error {
	err := c.storage.DB().WithContext(ctx).Where("id = ?", commentID).Delete(&model.Comment{}).Error
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return nil
}

func (c *commentRepository) GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).Model(&model.Comment{}).Where("film_id = ?", filmID).Count(&totalRecords).Error
	if err != nil {
		return []model.Comment{}, 0, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	err = c.storage.DB().WithContext(ctx).Where("film_id = ?", filmID).Limit(pageSize).Offset(offset).Find(&comments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Comment{}, 0, nil
	} else if err != nil {
		return []model.Comment{}, 0, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comments, totalRecords, nil
}

func (c *commentRepository) Update(ctx context.Context, commentDTO dto.UpdateCommentDTO, commentID string) (model.Comment, error) {
	var comment model.Comment
	err := c.storage.DB().WithContext(ctx).Model(&comment).Where("id = ?", commentID).Updates(commentDTO).Error
	if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	err = c.storage.DB().WithContext(ctx).First(&comment, "id = ?", commentID).Error
	if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comment, nil
}
