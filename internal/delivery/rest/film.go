package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FilmService interface {
	One(ctx context.Context, filmId string) (interface{}, error) // domain model, error
}

type FilmHandler struct {
	filmService FilmService
}

func NewFilmHandler(
	filmService FilmService,
) *FilmHandler {
	return &FilmHandler{
		filmService: filmService,
	}
}

func (h *FilmHandler) GetOneByID(ctx *fiber.Ctx) error {
	id := ctx.Get("id")
	film, err := h.filmService.One(ctx.Context(), id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		resp := &fiber.Map{"status": "NotFound", "error": errors.New("Film not found")}
		return ctx.JSON(resp)
	}

	return ctx.JSON(film)
}
