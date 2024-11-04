package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/pkg/token"
	"kbox-api/shared/httperror"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	jwtKey := c.Locals("jwtKey").(string)
	if tokenString == "" {
		return httperror.New(
			http.StatusUnauthorized,
			"Отсутствует токен",
		)
	}

	claims, err := token.ParseToken(tokenString, []byte(jwtKey))

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
