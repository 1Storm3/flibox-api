package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

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
	filmId := ctx.Params("id")
	film, err := h.filmService.GetOne(filmId)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		resp := &fiber.Map{"status": http.StatusInternalServerError, "error": errors.New("error in server")}
		return ctx.JSON(resp)
	}

	return ctx.JSON(film)
}
