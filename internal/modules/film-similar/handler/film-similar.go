package handler

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/modules/film-similar/service"
)

type FilmSimilarHandlerInterface interface {
	GetAll(ctx *fiber.Ctx) error
}

type FilmSimilarHandler struct {
	filmSimilarService service.FilmSimilarServiceInterface
}

func NewFilmSimilarHandler(
	filmSimilarService service.FilmSimilarServiceInterface,
) FilmSimilarHandlerInterface {
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
