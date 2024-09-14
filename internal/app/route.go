package app

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/external"
)

type routeConfig struct {
	app    *fiber.App
	config *config.Config
}

func newRouteConfig(app *fiber.App, config *config.Config) *routeConfig {
	return &routeConfig{
		app:    app,
		config: config,
	}
}

func (r *routeConfig) Setup() {
	r.app.Get("api/film/:id", func(ctx *fiber.Ctx) error {
		return external.GetFilm(ctx, r.config)
	})
}
