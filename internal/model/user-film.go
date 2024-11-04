package model

type UserFilm struct {
	UserId string `json:"userId" gorm:"column:user_id"`
	FilmId int    `json:"filmId" gorm:"column:film_id"`
	Film   Film   `gorm:"foreignKey:FilmId;references:ID"`
}
