package interceptor

import (
	"github.com/1Storm3/flibox-api/internal/metrics"
	"github.com/gofiber/fiber/v2"
)

func MetricsInterceptor() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		metrics.IncRequestCounter()

		return ctx.Next()
	}
}
