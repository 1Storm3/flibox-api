package model

type UserFilm struct {
	UserID string       `json:"userId" gorm:"column:user_id"`
	FilmID int          `json:"filmId" gorm:"column:film_id"`
	Type   TypeUserFilm `json:"type" gorm:"column:type"`
	Film   Film         `gorm:"foreignKey:FilmID;references:ID"`
}

type TypeUserFilm string

const (
	TypeUserFavourite TypeUserFilm = "favourite"
	TypeUserRecommend TypeUserFilm = "recommend"
)
