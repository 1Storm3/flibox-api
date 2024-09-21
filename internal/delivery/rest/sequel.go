package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/shared/logger"
	"net/http"
)

type SequelHandler struct {
	sequelService SequelService
}

func NewSequelHandler(
	sequelService SequelService,
) *SequelHandler {
	return &SequelHandler{
		sequelService: sequelService,
	}
}

func (h *SequelHandler) GetAll(ctx *fiber.Ctx) error {
	filmId := ctx.Params("id")
	sequels, err := h.sequelService.GetAll(filmId)
	if err != nil {
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		resp := fiber.Map{
			"error":      err.Error(),
			"statusCode": http.StatusInternalServerError,
		}
		return ctx.JSON(resp)
	}
	return ctx.JSON(sequels)
}
