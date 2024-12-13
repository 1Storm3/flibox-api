package collection

import (
	"kbox-api/internal/shared/httperror"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/pkg/token"
)

type HandlerInterface interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetOne(c *fiber.Ctx) error
	GetAllMy(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewCollectionHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Update(c *fiber.Ctx) error {
	collectionId := c.Params("id")
	ctx := c.Context()
	var collection UpdateCollectionDTO
	if err := c.BodyParser(&collection); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	result, err := h.service.Update(ctx, collection, collectionId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	collectionId := c.Params("id")
	ctx := c.Context()
	err := h.service.Delete(ctx, collectionId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": "Коллекция удалена",
	})
}

func (h *Handler) GetAllMy(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	result, totalRecords, err := h.service.GetAllMy(ctx, page, pageSize, userId)
	if err != nil {
		return err
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"data":         result,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	result, totalRecords, err := h.service.GetAll(ctx, page, pageSize)
	if err != nil {
		return err
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"data":         result,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *Handler) GetOne(c *fiber.Ctx) error {
	collectionId := c.Params("id")
	ctx := c.Context()
	result, err := h.service.GetOne(ctx, collectionId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID

	ctx := c.Context()
	var collection CreateCollectionDTO
	if err := c.BodyParser(&collection); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	result, err := h.service.Create(ctx, collection, userId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}
