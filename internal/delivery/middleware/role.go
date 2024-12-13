package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/shared/httperror"
	"kbox-api/pkg/token"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims := c.Locals("userClaims").(*token.Claims)
		userRole := userClaims.Role

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для выполнения операции",
		)
	}
}
