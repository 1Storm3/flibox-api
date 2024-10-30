package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/modules/film/handler"
	"kinopoisk-api/internal/modules/user-film/repository"
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
		logger.Error(err.Error())
		if isExist.Description == nil {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error":      "Фильм не найден",
				"statusCode": http.StatusNotFound,
			})
		}
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = g.userFilmService.Add(userId, filmId)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyAdded) {
			return ctx.JSON(fiber.Map{
				"message": "Фильм уже в избранном",
			})
		}
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"message": "Фильм добавлен в избранное",
	})
}

func (g *UserFilmHandler) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	filmId := ctx.Params("film_id")
	err := g.userFilmService.Delete(userId, filmId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ctx.JSON(fiber.Map{
				"message": "Фильм не в избранном",
			})
		}
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(err)
	}
	return ctx.JSON(fiber.Map{
		"message": "Фильм удален из избранного",
	})
}
