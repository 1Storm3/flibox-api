package auth

import (
	"errors"
	"kbox-api/internal/shared/httperror"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/pkg/token"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
}

type Handler struct {
	service ServiceInterface
}

func NewAuthHandler(authService ServiceInterface) *Handler {
	return &Handler{
		service: authService,
	}
}

// Login @Summary Login user
// @Description Login to the application
// @Accept  json
// @Produce  json
// @Param login body dto.LoginDTO true "Login information"
// @Success 200 token string
// @Failure 400 {object} httperror.Error
// @Router /auth/login [post]
func (a *Handler) Login(c *fiber.Ctx) error {
	var loginData LoginDTO

	ctx := c.Context()

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "Некорректные данные запроса",
			"statusCode": http.StatusBadRequest,
		})
	}
	tokenUser, err := a.service.Login(ctx, loginData)
	if err != nil {
		var httpErr *httperror.Error
		if errors.As(err, &httpErr) {
			return c.Status(httpErr.Code()).JSON(fiber.Map{
				"message":    httpErr.Error(),
				"statusCode": httpErr.Code(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    err.Error(),
			"statusCode": http.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenUser,
	})
}

func (a *Handler) Register(c *fiber.Ctx) error {
	var requestUser RegisterDTO

	ctx := c.Context()
	if err := c.BodyParser(&requestUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "Некорректные данные запроса",
			"statusCode": http.StatusBadRequest,
		})
	}
	result, err := a.service.Register(ctx, requestUser)
	if err != nil {
		var httpErr *httperror.Error
		if errors.As(err, &httpErr) {
			return c.Status(httpErr.Code()).JSON(fiber.Map{
				"message":    httpErr.Error(),
				"statusCode": httpErr.Code(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (a *Handler) Verify(c *fiber.Ctx) error {
	tokenUser := c.Params("token")
	if err := a.service.Verify(c.Context(), tokenUser); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Пользователь верифицирован",
	})
}

func (a *Handler) Me(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*token.Claims)

	ctx := c.Context()

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "Не удалось получить информацию о пользователе",
			"statusCode": http.StatusUnauthorized,
		})
	}
	result, err := a.service.Me(ctx, claims.UserID)

	if err != nil {
		var httpErr *httperror.Error
		if errors.As(err, &httpErr) {
			return c.Status(httpErr.Code()).JSON(fiber.Map{
				"message":    httpErr.Error(),
				"statusCode": httpErr.Code(),
			})
		}
	}

	return c.JSON(result)
}
