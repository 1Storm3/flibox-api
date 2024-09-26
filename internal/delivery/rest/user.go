package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/service"
	"kinopoisk-api/shared/logger"
	"net/http"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetOne(ctx *fiber.Ctx) error {
	userToken := ctx.Params("user_token")

	user, err := h.userService.GetOne(userToken)
	if err != nil {
		logger.Error(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if user == (service.User{}) {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":      "Пользователь с таким токеном не найден",
			"statusCode": http.StatusNotFound,
		})
	}
	return ctx.JSON(user)
}
