package external

import (
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/modules/external/handler"
	externalservice "kinopoisk-api/internal/modules/external/service"
)

type Module struct {
	config          *config.Config
	externalService handler.ExternalService
}

func NewExternalModule(config *config.Config) *Module {
	return &Module{
		config: config,
	}
}

func (m *Module) ExternalService() (handler.ExternalService, error) {
	if m.externalService == nil {
		m.externalService = externalservice.NewExternalService(m.config)
	}
	return m.externalService, nil
}
