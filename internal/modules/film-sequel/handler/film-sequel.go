package handler

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/modules/film-sequel/service"
)

type FilmSequelHandlerInterface interface {
	GetAll(ctx *fiber.Ctx) error
}

type FilmSequelHandler struct {
	filmSequelService service.FilmSequelServiceInterface
}

func NewFilmSequelHandler(
	filmSequelService service.FilmSequelServiceInterface,
) FilmSequelHandlerInterface {
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
