package userfilm

import (
	"context"
	"errors"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/model"
	filmService "github.com/1Storm3/flibox-api/internal/modules/film"
	"github.com/1Storm3/flibox-api/internal/modules/recommendation/adapter"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"github.com/1Storm3/flibox-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	GetAll(ctx *fiber.Ctx) error
	Add(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type Handler struct {
	service          ServiceInterface
	filmService      filmService.ServiceInterface
	recommendService adapter.RecommendService
}

func NewUserFilmHandler(
	service ServiceInterface,
	filmService filmService.ServiceInterface,
	recommendService adapter.RecommendService,
) *Handler {
	return &Handler{
		service:          service,
		filmService:      filmService,
		recommendService: recommendService,
	}
}

func (g *Handler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userClaims").(*token.Claims).UserID
	typeUserFilm := c.Query("type")

	ctx := c.Context()

	films, err := g.service.GetAll(ctx, userID, model.TypeUserFilm(typeUserFilm), 20)
	if err != nil {
		return err
	}

	return c.JSON(films)
}

func (g *Handler) Add(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := extractUserFilmParams(c)
	if err != nil {
		return err
	}

	ctx := c.Context()
	if err := g.checkFilmExistence(ctx, filmID); err != nil {
		return err
	}

	err = g.service.Add(ctx, Params{
		UserID: userID,
		FilmID: filmID,
		Type:   typeUserFilm,
	})
	if err != nil {
		return err
	}

	if typeUserFilm == model.TypeUserFavourite {
		go func() {
			err := g.recommendService.CreateRecommendations(adapter.RecommendationsParams{
				UserID: userID,
			})
			if err != nil {
				logger.Info("Произошла ошибка при создании рекомендаций")
				logger.Error(err.Error())
			}
		}()
	}

	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в избранное",
	})
}

func (g *Handler) Delete(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := extractUserFilmParams(c)
	if err != nil {
		return err
	}

	ctx := c.Context()
	err = g.service.Delete(ctx, Params{
		UserID: userID,
		FilmID: filmID,
		Type:   typeUserFilm,
	})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": "Фильм удален из избранного",
	})
}

func extractUserFilmParams(c *fiber.Ctx) (userID, filmID string, typeUserFilm model.TypeUserFilm, err error) {
	userID = c.Locals("userClaims").(*token.Claims).UserID
	filmID = c.Params("filmId")

	typeUserFilmReq := c.Query("type")
	if err := ParseTypeUserFilm(typeUserFilmReq, &typeUserFilm); err != nil {
		return "", "", "", httperror.New(
			http.StatusBadRequest,
			"Недопустимый тип фильма: "+typeUserFilmReq,
		)
	}
	return userID, filmID, typeUserFilm, nil
}

func (g *Handler) checkFilmExistence(ctx context.Context, filmID string) error {
	isExist, err := g.filmService.GetOne(ctx, filmID)
	if err != nil {
		return err
	}
	if isExist.Description == nil {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден",
		)
	}
	return nil
}

func ParseTypeUserFilm(s string, t *model.TypeUserFilm) error {
	switch s {
	case string(model.TypeUserFavourite):
		*t = model.TypeUserFavourite
	case string(model.TypeUserRecommend):
		*t = model.TypeUserRecommend
	default:
		return errors.New("Неверный тип фильма: " + s)
	}
	return nil
}
