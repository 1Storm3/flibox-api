package model

type FilmSimilar struct {
	SimilarId int  `json:"similarId" gorm:"column:similar_id"`
	FilmId    int  `json:"filmId" gorm:"column:film_id"`
	Film      Film `gorm:"foreignKey:FilmId;references:ID"`
}
