package model

type FilmSimilar struct {
	SimilarID int  `json:"similarId" gorm:"column:similar_id"`
	FilmID    int  `json:"filmId" gorm:"column:film_id"`
	Film      Film `gorm:"foreignKey:FilmId;references:ID"`
}
