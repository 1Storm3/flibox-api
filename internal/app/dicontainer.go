package app

import (
	"kinopoisk-api/internal/config"
)

type diContainer struct {
	config *config.Config
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
