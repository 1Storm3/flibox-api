package filmsequel

import (
	"github.com/gofiber/fiber/v2"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	GetAll(ctx *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewFilmSequelHandler(
	service ServiceInterface,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	filmId := c.Params("id")
	ctx := c.Context()
	sequels, err := h.service.GetAll(ctx, filmId)
	if err != nil {
		return err
	}
	return c.JSON(sequels)
}
