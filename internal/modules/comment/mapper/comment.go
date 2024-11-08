package mapper

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/comment/dto"
)

func MapModelCommentToResponseDTO(comment model.Comment) dto.CommentResponseDTO {
	return dto.CommentResponseDTO{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		FilmID:    comment.FilmID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt.String(),
		UpdatedAt: comment.UpdatedAt.String(),
	}
}
