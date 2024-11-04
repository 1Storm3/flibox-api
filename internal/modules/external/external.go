package external

import (
	"kbox-api/internal/config"
	"kbox-api/internal/modules/external/handler"
	"kbox-api/internal/modules/external/service"
)

type Module struct {
	config          *config.Config
	externalService handler.ExternalService
	s3Service       *service.S3Service
	externalHandler *handler.ExternalHandler
}

func NewExternalModule(config *config.Config) *Module {
	return &Module{
		config: config,
	}
}

func (m *Module) ExternalService() (handler.ExternalService, error) {
	if m.externalService == nil {
		m.externalService = service.NewExternalService(m.config)
	}
	return m.externalService, nil
}

func (m *Module) S3Service() (*service.S3Service, error) {
	if m.s3Service == nil {

		s3Service, err := service.NewS3Service(m.config)
		if err != nil {
			return nil, err
		}
		m.s3Service = s3Service
	}
	return m.s3Service, nil
}

func (m *Module) ExternalHandler() (*handler.ExternalHandler, error) {
	if m.externalHandler == nil && m.s3Service == nil {
		externalService, err := m.ExternalService()
		if err != nil {
			return nil, err
		}
		s3Service, err := m.S3Service()
		if err != nil {
			return nil, err
		}

		m.externalHandler = handler.NewExternalHandler(externalService, s3Service)
	}
	return m.externalHandler, nil
}
