package model

type CollectionFilm struct {
	CollectionID string     `json:"collectionId" gorm:"column:collection_id"`
	FilmID       int        `json:"filmId" gorm:"column:film_id"`
	Film         Film       `gorm:"foreignKey:FilmID;references:ID"`
	Collection   Collection `gorm:"foreignKey:CollectionID;references:ID"`
}
