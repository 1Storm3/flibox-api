package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
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
		return err
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
		return err
	}

	return ctx.JSON(film)
}
