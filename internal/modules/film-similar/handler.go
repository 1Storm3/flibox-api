package filmsimilar

import (
	"github.com/gofiber/fiber/v2"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	GetAll(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewFilmSimilarHandler(
	service ServiceInterface,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	filmId := c.Params("id")

	ctx := c.Context()

	similars, err := h.service.GetAll(ctx, filmId)
	if err != nil {
		return err
	}
	return c.JSON(similars)
}
