package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/shared/logger"
	"net/http"
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
		logger.Error(err.Error())
		ctx.Status(http.StatusConflict)
		resp := fiber.Map{
			"error":      err.Error(),
			"statusCode": http.StatusConflict,
		}
		return ctx.JSON(resp)
	}
	return ctx.JSON(similars)
}
