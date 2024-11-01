package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/pkg/token"
	"kinopoisk-api/shared/httperror"
	"net/http"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return httperror.New(
			http.StatusUnauthorized,
			"Отсутствует токен",
		)
	}

	claims, err := token.ParseToken(tokenString)

	if err != nil {
		_ = fmt.Errorf("Ошибка при разборе токена: %w", err)
		return httperror.New(
			http.StatusUnauthorized,
			"Недействительный токен",
		)
	}

	c.Locals("userClaims", claims)
	return c.Next()
}
