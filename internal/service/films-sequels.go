package service

type FilmsSequel struct {
	SequelId int `json:"sequelId" gorm:"column:sequel_id"`
	FilmId   int `json:"filmId" gorm:"column:film_id"`
}
