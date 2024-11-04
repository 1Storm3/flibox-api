package model

type FilmSequel struct {
	SequelId int  `json:"sequelId" gorm:"column:sequel_id"`
	FilmId   int  `json:"filmId" gorm:"column:film_id"`
	Film     Film `gorm:"foreignKey:FilmId;references:ID"`
}
