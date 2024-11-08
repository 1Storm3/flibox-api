package model

type FilmSequel struct {
	SequelID int  `json:"sequelId" gorm:"column:sequel_id"`
	FilmID   int  `json:"filmId" gorm:"column:film_id"`
	Film     Film `gorm:"foreignKey:FilmId;references:ID"`
}
