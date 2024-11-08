package model

import "time"

type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	FilmID    int       `json:"filmId" gorm:"column:film_id"`
	UserID    string    `json:"userId" gorm:"column:user_id"`
	ParentID  *string   `json:"parentId" gorm:"column:parent_id"`
	Content   string    `json:"content" gorm:"column:content"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
