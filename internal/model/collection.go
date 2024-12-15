package model

import (
	"github.com/lib/pq"

	"time"
)

type Collection struct {
	ID          string          `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string          `json:"name" gorm:"column:name"`
	Description *string         `json:"description" gorm:"column:description"`
	CoverUrl    *string         `json:"coverUrl" gorm:"column:cover_url"`
	Tags        *pq.StringArray `json:"tags" gorm:"type:text[];column:tags"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	UserId      *string         `json:"userId" gorm:"column:user_id"`
	User        User            `gorm:"foreignKey:UserId;references:ID"`
	Films       []Film          `gorm:"many2many:collection_films;joinForeignKey:collection_id;JoinReferences:film_id"`
}
