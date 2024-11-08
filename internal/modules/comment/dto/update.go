package dto

type UpdateCommentDTO struct {
	Content string `json:"content" validate:"required,min=1"`
}
