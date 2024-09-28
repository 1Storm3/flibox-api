package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/shared/logger"
	"net/http"
)

type SequelHandler struct {
	filmSequelService FilmSequelService
}

func NewSequelHandler(
	filmSequelService FilmSequelService,
) *SequelHandler {
	return &SequelHandler{
		filmSequelService: filmSequelService,
	}
}

func (h *SequelHandler) GetAll(ctx *fiber.Ctx) error {
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
