package model

import "time"

type HistoryFilms struct {
	ID        string    `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID    string    `json:"userId" gorm:"column:user_id"`
	FilmID    int       `json:"filmId" gorm:"column:film_id"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Film      Film      `gorm:"foreignKey:FilmID"`
	User      User      `gorm:"foreignKey:UserID"`
}
