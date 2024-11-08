package model

type UserFilm struct {
	UserID string `json:"userId" gorm:"column:user_id"`
	FilmID int    `json:"filmId" gorm:"column:film_id"`
	Film   Film   `gorm:"foreignKey:FilmId;references:ID"`
}
