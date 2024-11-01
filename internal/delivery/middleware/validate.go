package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateMiddleware[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqBody := new(T)
		if err := c.BodyParser(reqBody); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"statusCode": http.StatusBadRequest,
				"message":    "Некорректные данные запроса",
			})
		}

		if err := validate.Struct(reqBody); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				var errorMessages []string
				for _, fe := range ve {
					switch fe.Tag() {
					case "min":
						errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' должно содержать минимум %s символов", fe.Field(), fe.Param()))
					case "required":
						errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' обязательно для заполнения", fe.Field()))
					case "email":
						errorMessages = append(errorMessages, "Некорректный формат email")
					default:
						errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' некорректно", fe.Field()))
					}

				}
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"statusCode": http.StatusBadRequest,
					"message":    errorMessages[0],
				})
			}
		}
		return c.Next()
	}
}
