package external

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"kbox-api/internal/shared/httperror"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	UploadFile(c *fiber.Ctx) error
}

type Handler struct {
	service   ServiceInterface
	s3Service S3ServiceInterface
}

func NewExternalHandler(
	service ServiceInterface,
	s3Service S3ServiceInterface,
) *Handler {
	return &Handler{
		service:   service,
		s3Service: s3Service,
	}
}

func (e *Handler) UploadFile(c *fiber.Ctx) error {
	ctx := c.Context()
	file, err := c.FormFile("file")
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			"Не удалось получить данные",
		)
	}
	fileReader, err := file.Open()
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			"Не удалось получить данные",
		)
	}
	defer func(fileReader multipart.File) {
		err := fileReader.Close()
		if err != nil {
			return
		}
	}(fileReader)

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			"Не удалось прочитать файл",
		)
	}

	ext := filepath.Ext(file.Filename)
	uniqueID, err := uuid.NewUUID()
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			"Не удалось создать уникальный идентификатор",
		)
	}

	uniqueFilename := fmt.Sprintf("%s%s", uniqueID.String(), ext)

	url, err := e.s3Service.UploadFile(ctx, uniqueFilename, fileBytes)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}
