package handler

import (
	"kbox-api/internal/modules/user-film/dto"
)

type UserFilmService interface {
	GetAll(userId string) ([]dto.GetUserFilmResponseDTO, error)
	Add(userId string, filmId string) error
	Delete(userId string, filmId string) error
}
