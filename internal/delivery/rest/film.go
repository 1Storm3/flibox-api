package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/shared/httperror"
	"kinopoisk-api/shared/logger"
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
		logger.Error(err.Error())

		var httpErr *httperror.Error
		if errors.As(err, &httpErr) && httpErr.Code() == http.StatusNotFound {
			ctx.Status(http.StatusNotFound)
			resp := fiber.Map{
				"error":      httpErr.Error(),
				"statusCode": http.StatusNotFound,
			}
			return ctx.JSON(resp)
		}

		ctx.Status(http.StatusInternalServerError)
		resp := fiber.Map{
			"error": httpErr.Error(),
			"code":  http.StatusInternalServerError,
		}
		return ctx.JSON(resp)
	}

	return ctx.JSON(film)
}
