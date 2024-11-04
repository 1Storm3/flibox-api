package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"kbox-api/internal/modules/user/dto"
	"kbox-api/internal/modules/user/mapper"
	"kbox-api/shared/httperror"
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

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var userUpdateRequest dto.UpdateUserDTO
	if err := c.BodyParser(&userUpdateRequest); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	userUpdateRequest.Id = id

	result, err := h.userService.Update(userUpdateRequest)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
		}
		return err
	}
	updatedUser := mapper.MapModelUserToResponseDTO(result)

	return c.JSON(updatedUser)
}
