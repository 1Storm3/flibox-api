package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/shared/logger"
	"net/http"
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
		logger.Error(err.Error())
		ctx.Status(http.StatusConflict)
		resp := fiber.Map{
			"error":      err.Error(),
			"statusCode": http.StatusConflict,
		}
		return ctx.JSON(resp)
	}
	return ctx.JSON(sequels)
}
