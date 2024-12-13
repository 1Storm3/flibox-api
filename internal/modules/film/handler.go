package film

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	Search(c *fiber.Ctx) error
	GetOneByID(c *fiber.Ctx) error
}
type Handler struct {
	service ServiceInterface
}

func NewFilmHandler(
	service ServiceInterface,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Search(c *fiber.Ctx) error {
	match := c.Query("match")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	genresStr := c.Query("genres")
	ctx := c.Context()
	var genres []string
	if genresStr != "" {
		genres = strings.Split(genresStr, ",")
	}
	films, totalRecords, err := h.service.Search(ctx, match, genres, page, pageSize)

	if err != nil {
		return err
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"films":        films,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}
func (h *Handler) GetOneByID(c *fiber.Ctx) error {
	filmId := c.Params("id")
	ctx := c.Context()
	film, err := h.service.GetOne(ctx, filmId)

	if err != nil {
		return err
	}

	return c.JSON(film)
}
