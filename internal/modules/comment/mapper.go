package comment

import "github.com/1Storm3/flibox-api/internal/model"

func MapModelCommentToResponseDTO(comment model.Comment) ResponseDTO {
	return ResponseDTO{
		ID:      comment.ID,
		Content: comment.Content,
		User: User{
			ID:       comment.User.ID,
			NickName: comment.User.NickName,
			Photo:    comment.User.Photo,
		},
		FilmID:    comment.FilmID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
