package collectionfilm

import (
	"net/http"

	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
)

type HandlerInterface interface {
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetFilmsByCollectionId(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewCollectionFilmHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Add(c *fiber.Ctx) error {
	var filmId CreateCollectionFilmDTO
	collectionId := c.Params("id")
	if err := c.BodyParser(&filmId); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	ctx := c.Context()
	err := h.service.Add(ctx, collectionId, filmId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в подборку",
	})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	var filmId DeleteCollectionFilmDTO
	collectionId := c.Params("id")
	if err := c.BodyParser(&filmId); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	ctx := c.Context()
	err := h.service.Delete(ctx, collectionId, filmId)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": "Фильмы удалены из подборки",
	})
}

func (h *Handler) GetFilmsByCollectionId(c *fiber.Ctx) error {
	collectionID := c.Params("id")
	if collectionID == "" {
		return httperror.New(
			http.StatusBadRequest,
			"Неверный формат ID коллекции",
		)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	ctx := c.Context()
	films, totalRecords, err := h.service.GetFilmsByCollectionId(ctx, collectionID, page, pageSize)
	if err != nil {
		return err
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"data":         films,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}
