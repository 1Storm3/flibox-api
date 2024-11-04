package handler

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/auth/dto"
)

type AuthService interface {
	Login(dto dto.LoginDTO) (string, error)
	Register(user dto.RegisterDTO) (string, error)
	Me(c *fiber.Ctx) (model.User, error)
}
