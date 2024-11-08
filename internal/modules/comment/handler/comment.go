package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/modules/comment/dto"
	"kbox-api/internal/modules/comment/service"
	"kbox-api/pkg/token"
	"kbox-api/shared/httperror"
)

var _ CommentHandlerInterface = (*CommentHandler)(nil)

type CommentHandlerInterface interface {
	GetAllByFilmID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type CommentHandler struct {
	commentService service.CommentServiceInterface
}

func NewCommentHandler(commentService service.CommentServiceInterface) CommentHandlerInterface {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) GetAllByFilmID(c *fiber.Ctx) error {
	filmID := c.Params("filmId")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	comments, totalRecords, err := h.commentService.GetAllByFilmId(filmID, page, pageSize)
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

func (h *CommentHandler) Create(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID

	var comment dto.CreateCommentDTO
	if err := c.BodyParser(&comment); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	result, err := h.commentService.Create(comment, userId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *CommentHandler) Update(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID
	role := c.Locals("userClaims").(*token.Claims).Role

	commentId := c.Params("id")
	comment, err := h.commentService.GetOne(commentId)
	if err != nil {
		return err
	}

	if role != "admin" && comment.UserID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для редактирования комментария",
		)
	}
	var commentDto dto.UpdateCommentDTO
	if err := c.BodyParser(&commentDto); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	result, err := h.commentService.Update(commentDto, commentId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*token.Claims).UserID
	role := c.Locals("userClaims").(*token.Claims).Role
	commentId := c.Params("id")

	comment, err := h.commentService.GetOne(commentId)
	if err != nil {
		return err
	}

	if role != "admin" && comment.UserID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для удаления комментария",
		)
	}
	err = h.commentService.Delete(commentId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": "Комментарий удален",
	})
}
