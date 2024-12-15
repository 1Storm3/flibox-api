package user

import (
	"errors"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	GetOneByNickName(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewUserHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetOneByNickName(c *fiber.Ctx) error {
	nickName := c.Params("nickName")

	ctx := c.Context()

	user, err := h.service.GetOneByNickName(ctx, nickName)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx := c.Context()

	var userUpdateRequest UpdateUserDTO
	if err := c.BodyParser(&userUpdateRequest); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	userUpdateRequest.ID = id

	result, err := h.service.Update(ctx, userUpdateRequest)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
		}
		return err
	}
	updatedUser := MapModelUserToResponseDTO(result)

	return c.JSON(updatedUser)
}
