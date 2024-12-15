package comment

import (
	"net/http"
	"strings"

	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	GetAllByFilmID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewCommentHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAllByFilmID(c *fiber.Ctx) error {
	filmID := c.Params("filmId")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	comments, totalRecords, err := h.service.GetAllByFilmId(ctx, filmID, page, pageSize)
	if err != nil {
		return err
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"comments":     comments,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID

	ctx := c.Context()
	var comment CreateCommentDTO
	if err := c.BodyParser(&comment); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	if len(strings.TrimSpace(*comment.Content)) == 0 {
		return httperror.New(http.StatusBadRequest, "Комментарий не может быть пустым")
	}
	result, err := h.service.Create(ctx, comment, userId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID
	role := c.Locals("userClaims").(*token.Claims).Role

	commentId := c.Params("id")

	ctx := c.Context()
	comment, err := h.service.GetOne(ctx, commentId)
	if err != nil {
		return err
	}

	if role != "admin" && comment.User.ID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для редактирования комментария",
		)
	}
	var commentDto UpdateCommentDTO
	if err := c.BodyParser(&commentDto); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	result, err := h.service.Update(ctx, commentDto, commentId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID
	role := c.Locals("userClaims").(*token.Claims).Role
	commentId := c.Params("id")
	ctx := c.Context()
	comment, err := h.service.GetOne(ctx, commentId)
	if err != nil {
		return err
	}

	if role != "admin" && comment.User.ID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для удаления комментария",
		)
	}
	err = h.service.Delete(ctx, commentId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": "Комментарий удален",
	})
}
