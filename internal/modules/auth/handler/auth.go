package handler

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/modules/auth/service"
	"kinopoisk-api/shared/httperror"
	"net/http"
)

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	var loginData RequestLogin
	if err := c.BodyParser(&loginData); err != nil {

		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	token, err := a.authService.Login(loginData.Email, loginData.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (a *AuthHandler) Register(c *fiber.Ctx) error {
	var requestUser service.RequestUser
	if err := c.BodyParser(&requestUser); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	tokenString, err := a.authService.Register(requestUser)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
func (a *AuthHandler) Me(c *fiber.Ctx) error {
	user, err := a.authService.Me(c)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
