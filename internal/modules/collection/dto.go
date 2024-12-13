package collection

import (
	"github.com/lib/pq"
	"time"
)

type CreateCollectionDTO struct {
	Name        string          `json:"name" validate:"required,min=4"`
	Description *string         `json:"description" validate:"omitempty,min=4"`
	CoverUrl    *string         `json:"coverUrl" validate:"omitempty,url"`
	Tags        *pq.StringArray `json:"tags" validate:"omitempty,min=1"`
}

type ResponseDTO struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	CoverUrl    *string         `json:"coverUrl"`
	User        User            `json:"author"`
	Tags        *pq.StringArray `json:"tags"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type UpdateCollectionDTO struct {
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	CoverUrl    *string         `json:"coverUrl"`
	Tags        *pq.StringArray `json:"tags"`
}

type Film struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	PosterURL       *string  `json:"posterUrl"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk"`
}

type User struct {
	ID       string  `json:"id"`
	NickName string  `json:"nickName"`
	Photo    *string `json:"photo"`
}
