package interceptor

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/metrics"
)

func MetricsInterceptor() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		metrics.IncRequestCounter()

		return ctx.Next()
	}
}
