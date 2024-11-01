package handler

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetOneByNickName(c *fiber.Ctx) error {
	nickName := c.Params("nickName")

	user, err := h.userService.GetOneByNickName(nickName)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
