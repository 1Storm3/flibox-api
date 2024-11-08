package service

import (
	"context"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/comment/dto"
	"kbox-api/internal/modules/comment/mapper"
	"kbox-api/internal/modules/comment/repository"
)

var _ CommentServiceInterface = (*CommentService)(nil)

type CommentServiceInterface interface {
	Create(comment dto.CreateCommentDTO, userID string) (dto.CommentResponseDTO, error)
	Update(comment dto.UpdateCommentDTO, commentID string) (dto.CommentResponseDTO, error)
	Delete(commentID string) error
	GetAllByFilmId(filmID string, page int, pageSize int) ([]dto.CommentResponseDTO, int64, error)
	GetOne(commentID string) (dto.CommentResponseDTO, error)
}

type CommentService struct {
	commentRepo repository.CommentRepositoryInterface
}

func NewCommentService(commentRepo repository.CommentRepositoryInterface) CommentServiceInterface {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (c *CommentService) Create(comment dto.CreateCommentDTO, userID string) (dto.CommentResponseDTO, error) {
	result, err := c.commentRepo.Create(context.Background(), model.Comment{
		Content:  comment.Content,
		FilmID:   comment.FilmID,
		UserID:   userID,
		ParentID: comment.ParentID,
	})

	if err != nil {
		return dto.CommentResponseDTO{}, err
	}
	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) GetOne(commentID string) (dto.CommentResponseDTO, error) {
	result, err := c.commentRepo.GetOne(context.Background(), commentID)
	if err != nil {
		return dto.CommentResponseDTO{}, err
	}

	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) Update(comment dto.UpdateCommentDTO, commentID string) (dto.CommentResponseDTO, error) {
	result, err := c.commentRepo.Update(context.Background(), comment, commentID)
	if err != nil {
		return dto.CommentResponseDTO{}, err
	}
	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) Delete(commentID string) error {
	return c.commentRepo.Delete(context.Background(), commentID)
}

func (c *CommentService) GetAllByFilmId(filmID string, page int, pageSize int) ([]dto.CommentResponseDTO, int64, error) {
	comments, totalRecords, err := c.commentRepo.GetAllByFilmId(context.Background(), filmID, page, pageSize)

	if err != nil {
		return []dto.CommentResponseDTO{}, 0, err
	}
	var commentsDTO []dto.CommentResponseDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, mapper.MapModelCommentToResponseDTO(comment))
	}
	if len(commentsDTO) == 0 {
		return []dto.CommentResponseDTO{}, totalRecords, nil
	}
	return commentsDTO, totalRecords, nil
}
