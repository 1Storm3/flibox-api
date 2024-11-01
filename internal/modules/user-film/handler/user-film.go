package handler

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/modules/film/handler"
	"kinopoisk-api/shared/httperror"
	"kinopoisk-api/shared/logger"
	"net/http"
)

type UserFilmHandler struct {
	userFilmService UserFilmService
	filmService     handler.FilmService
}

func NewUserFilmHandler(
	userFilmService UserFilmService,
	filmService handler.FilmService,
) *UserFilmHandler {
	return &UserFilmHandler{
		userFilmService: userFilmService,
		filmService:     filmService,
	}
}

func (g *UserFilmHandler) GetAll(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")

	films, err := g.userFilmService.GetAll(userId)
	if err != nil {
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(err)
	}
	return ctx.JSON(films)
}

func (g *UserFilmHandler) Add(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	filmId := ctx.Params("film_id")
	isExist, err := g.filmService.GetOne(filmId)
	if err != nil {
		return err
	}
	if isExist.Description == nil {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден",
		)
	}

	err = g.userFilmService.Add(userId, filmId)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"data": "Фильм добавлен в избранное",
	})
}

func (g *UserFilmHandler) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	filmId := ctx.Params("film_id")
	err := g.userFilmService.Delete(userId, filmId)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"data": "Фильм удален из избранного",
	})
}
