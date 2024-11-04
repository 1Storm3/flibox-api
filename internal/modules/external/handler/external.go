package handler

import (
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/shared/httperror"
)

type ExternalHandler struct {
	ExternalService ExternalService
	S3Service       S3Service
}

func NewExternalHandler(externalService ExternalService, s3Service S3Service) *ExternalHandler {
	return &ExternalHandler{
		ExternalService: externalService,
		S3Service:       s3Service,
	}
}

func (e *ExternalHandler) UploadFile(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
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
	defer fileReader.Close()

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			"Не удалось прочитать файл",
		)
	}

	url, err := e.S3Service.UploadFile(ctx.Context(), file.Filename, fileBytes)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}
