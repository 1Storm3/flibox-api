package handler

import (
	"github.com/gofiber/fiber/v2"
)

type FilmSequelHandler struct {
	filmSequelService FilmSequelService
}

func NewFilmSequelHandler(
	filmSequelService FilmSequelService,
) *FilmSequelHandler {
	return &FilmSequelHandler{
		filmSequelService: filmSequelService,
	}
}

func (h *FilmSequelHandler) GetAll(ctx *fiber.Ctx) error {
	filmId := ctx.Params("id")
	sequels, err := h.filmSequelService.GetAll(filmId)
	if err != nil {
		return err
	}
	return ctx.JSON(sequels)
}
