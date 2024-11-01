package handler

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/modules/auth/service"
	userservice "kinopoisk-api/internal/modules/user/service"
)

type AuthService interface {
	Login(email string, password string) (string, error)
	Register(user service.RequestUser) (string, error)
	Me(c *fiber.Ctx) (userservice.User, error)
}
