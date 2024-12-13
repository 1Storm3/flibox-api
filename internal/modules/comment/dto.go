package comment

import "time"

type CreateCommentDTO struct {
	Content  *string `json:"content" validate:"required,min=1"`
	FilmID   int     `json:"filmId" validate:"required,min=36"`
	ParentID *string `json:"parentId" validate:"omitempty"`
}

type ResponseDTO struct {
	ID        string    `json:"id"`
	Content   *string   `json:"content"`
	FilmID    int       `json:"filmId"`
	User      User      `json:"author"`
	ParentID  *string   `json:"parentId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateCommentDTO struct {
	Content *string `json:"content" validate:"required,min=1"`
}

type User struct {
	ID       string  `json:"id"`
	NickName string  `json:"nickName"`
	Photo    *string `json:"photo"`
}
