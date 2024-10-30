package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"kinopoisk-api/shared/httperror"
	"kinopoisk-api/shared/logger"
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

func (h *FilmHandler) Search(ctx *fiber.Ctx) error {
	match := ctx.Query("match")
	page := ctx.QueryInt("page", 1)
	pageSize := ctx.QueryInt("pageSize", 10)
	genresStr := ctx.Query("genres")
	var genres []string
	if genresStr != "" {
		genres = strings.Split(genresStr, ",")
	}
	films, totalRecords, err := h.filmService.Search(match, genres, page, pageSize)

	if err != nil {
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(err)
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"films":        films,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
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
			"error": "Internal Server Error",
			"code":  http.StatusInternalServerError,
		}
		return ctx.JSON(resp)
	}

	return ctx.JSON(film)
}
