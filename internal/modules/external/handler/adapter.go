package handler

import (
	"context"
	"kbox-api/internal/modules/external/dto"
)

type ExternalService interface {
	SetExternalRequest(filmId string) (dto.GetExternalFilmDTO, error)
}

type S3Service interface {
	UploadFile(ctx context.Context, key string, file []byte) (string, error)
}
