package comment

import (
	"context"
	"kbox-api/internal/model"
)

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	Create(ctx context.Context, comment CreateCommentDTO, userID string) (ResponseDTO, error)
	Update(ctx context.Context, comment UpdateCommentDTO, commentID string) (ResponseDTO, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page int, pageSize int) ([]ResponseDTO, int64, error)
	GetOne(ctx context.Context, commentID string) (ResponseDTO, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewCommentService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (c *Service) Create(ctx context.Context, comment CreateCommentDTO, userID string) (ResponseDTO, error) {
	result, err := c.repository.Create(ctx, model.Comment{
		Content:  comment.Content,
		FilmID:   comment.FilmID,
		UserID:   userID,
		ParentID: comment.ParentID,
	})

	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelCommentToResponseDTO(result), nil
}

func (c *Service) GetOne(ctx context.Context, commentID string) (ResponseDTO, error) {
	result, err := c.repository.GetOne(ctx, commentID)
	if err != nil {
		return ResponseDTO{}, err
	}

	return MapModelCommentToResponseDTO(result), nil
}

func (c *Service) Update(ctx context.Context, comment UpdateCommentDTO, commentID string) (ResponseDTO, error) {
	result, err := c.repository.Update(ctx, comment, commentID)
	if err != nil {
		return ResponseDTO{}, err
	}
	return MapModelCommentToResponseDTO(result), nil
}

func (c *Service) Delete(ctx context.Context, commentID string) error {
	comment, err := c.repository.GetOne(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.ParentID == nil {
		countChildComments, err := c.repository.GetCountByParentId(ctx, commentID)
		if err != nil {
			return err
		}
		if countChildComments != 0 {
			_, err := c.repository.Update(ctx, UpdateCommentDTO{Content: nil}, commentID)
			if err != nil {
				return err
			}
		} else {
			err := c.repository.Delete(ctx, commentID)
			if err != nil {
				return err
			}
		}
	} else {
		countSiblingComment, err := c.repository.GetCountByParentId(ctx, *comment.ParentID)
		if err != nil {
			return err
		}
		if countSiblingComment == 1 {
			parentComment, err := c.repository.GetOne(ctx, *comment.ParentID)
			if err != nil {
				return err
			}
			if parentComment.Content == nil {
				err := c.repository.Delete(ctx, *comment.ParentID)
				if err != nil {
					return err
				}
			}
		}

		err = c.repository.Delete(ctx, commentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Service) GetAllByFilmId(ctx context.Context, filmID string, page int, pageSize int) ([]ResponseDTO, int64, error) {
	comments, totalRecords, err := c.repository.GetAllByFilmId(ctx, filmID, page, pageSize)
	if err != nil {
		return []ResponseDTO{}, 0, err
	}
	var commentsDTO []ResponseDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, MapModelCommentToResponseDTO(comment))
	}
	if len(commentsDTO) == 0 {
		return []ResponseDTO{}, totalRecords, nil
	}
	return commentsDTO, totalRecords, nil
}
