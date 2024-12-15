package middleware

import (
	"net/http"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/modules/user"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(userRepo user.RepositoryInterface, config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		jwtKey := config.App.JwtSecretKey
		if tokenString == "" {
			return httperror.New(
				http.StatusUnauthorized,
				"Отсутствует токен",
			)
		}

		claims, err := token.ParseToken(tokenString, []byte(jwtKey))
		if err != nil {
			return httperror.New(
				http.StatusUnauthorized,
				"Недействительный токен")
		}

		user, err := userRepo.GetOneById(c.Context(), claims.UserID)
		if err != nil {
			return httperror.New(
				http.StatusUnauthorized,
				"Ошибка получения информации о пользователе",
			)
		}
		if user.IsBlocked {
			return httperror.New(
				http.StatusForbidden,
				"Пользователь заблокирован")
		}
		if !user.IsVerified {
			return httperror.New(
				http.StatusForbidden,
				"Пользователь не верифицирован")
		}

		c.Locals("userClaims", claims)
		return c.Next()
	}
}
