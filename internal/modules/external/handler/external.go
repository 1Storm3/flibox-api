package handler

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/modules/external/service"
	"kbox-api/shared/httperror"
)

type ExternalHandlerInterface interface {
	UploadFile(ctx *fiber.Ctx) error
}

type ExternalHandler struct {
	ExternalService service.ExternalServiceInterface
	S3Service       service.S3ServiceInterface
}

func NewExternalHandler(
	externalService service.ExternalServiceInterface,
	s3Service service.S3ServiceInterface,
) ExternalHandlerInterface {
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
	url, err := e.S3Service.UploadFile(ctx.Context(), file.Filename, fileBytes)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}
