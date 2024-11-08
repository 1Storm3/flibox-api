package dto

type CreateCommentDTO struct {
	Content  string  `json:"content" validate:"required,min=1"`
	FilmID   int     `json:"filmId" validate:"required,min=36"`
	ParentID *string `json:"parentId" validate:"omitempty"`
}

type CommentResponseDTO struct {
	ID        string  `json:"id"`
	Content   string  `json:"content"`
	FilmID    int     `json:"filmId"`
	UserID    string  `json:"userId"`
	ParentID  *string `json:"parentId,omitempty"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
