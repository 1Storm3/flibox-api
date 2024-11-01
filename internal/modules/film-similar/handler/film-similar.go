package handler

import (
	"github.com/gofiber/fiber/v2"
)

type FilmSimilarHandler struct {
	filmSimilarService FilmSimilarService
}

func NewFilmSimilarHandler(
	filmSimilarService FilmSimilarService,
) *FilmSimilarHandler {
	return &FilmSimilarHandler{
		filmSimilarService: filmSimilarService,
	}
}

func (h *FilmSimilarHandler) GetAll(ctx *fiber.Ctx) error {
	filmId := ctx.Params("id")
	similars, err := h.filmSimilarService.GetAll(filmId)
	if err != nil {
		return err
	}
	return ctx.JSON(similars)
}
