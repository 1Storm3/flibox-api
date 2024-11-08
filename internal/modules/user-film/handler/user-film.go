package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	filmService "kbox-api/internal/modules/film/service"
	"kbox-api/internal/modules/user-film/service"
	"kbox-api/pkg/token"
	"kbox-api/shared/httperror"
)

var _ UserFilmHandlerInterface = (*UserFilmHandler)(nil)

type UserFilmHandlerInterface interface {
	GetAll(ctx *fiber.Ctx) error
	Add(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type UserFilmHandler struct {
	userFilmService service.UserFilmServiceInterface
	filmService     filmService.FilmServiceInterface
}

func NewUserFilmHandler(
	userFilmService service.UserFilmServiceInterface,
	filmService filmService.FilmServiceInterface,
) *UserFilmHandler {
	return &UserFilmHandler{
		userFilmService: userFilmService,
		filmService:     filmService,
	}
}

func (g *UserFilmHandler) GetAll(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userClaims").(*token.Claims).UserID

	films, err := g.userFilmService.GetAll(userId)

	if err != nil {
		return err
	}
	return ctx.JSON(films)
}

func (g *UserFilmHandler) Add(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userClaims").(*token.Claims).UserID
	filmId := ctx.Params("filmId")
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
	userId := ctx.Locals("userClaims").(*token.Claims).UserID
	filmId := ctx.Params("filmId")
	err := g.userFilmService.Delete(userId, filmId)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"data": "Фильм удален из избранного",
	})
}
