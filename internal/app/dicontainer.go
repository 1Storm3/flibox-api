package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/config"
)

type diContainer struct {
	config *config.Config
	app    *fiber.App
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Config() *config.Config {
	if d.config == nil {
		d.config = config.MustLoad()
	}
	return d.config
}

func (d *diContainer) App(_ context.Context) (*fiber.App, error) {
	if d.app == nil {
		d.app = fiber.New()
	}

	return d.app, nil
}
