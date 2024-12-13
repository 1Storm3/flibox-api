package collectionfilm

import (
	"kbox-api/internal/modules/collection"
)

type CreateCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type DeleteCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type FilmsByCollectionIdDTO struct {
	CollectionID string            `json:"collectionId" validate:"required"`
	Films        []collection.Film `json:"films" validate:"required"`
}
