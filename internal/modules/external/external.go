package external

import (
	"kbox-api/internal/config"
	"kbox-api/internal/modules/external/handler"
	"kbox-api/internal/modules/external/service"
)

type ModuleInterface interface {
	ExternalService() (service.ExternalServiceInterface, error)
	S3Service() (service.S3ServiceInterface, error)
	ExternalHandler() (handler.ExternalHandlerInterface, error)
}

type Module struct {
	cfg             *config.Config
	externalService service.ExternalServiceInterface
	s3Service       service.S3ServiceInterface
	externalHandler handler.ExternalHandlerInterface
}

func NewExternalModule(cfg *config.Config) ModuleInterface {
	return &Module{
		cfg: cfg,
	}
}

func (m *Module) ExternalService() (service.ExternalServiceInterface, error) {
	if m.externalService == nil {
		m.externalService = service.NewExternalService(m.cfg)
	}
	return m.externalService, nil
}

func (m *Module) S3Service() (service.S3ServiceInterface, error) {
	if m.s3Service == nil {

		s3Service, err := service.NewS3Service(m.cfg)
		if err != nil {
			return nil, err
		}
		m.s3Service = s3Service
	}
	return m.s3Service, nil
}

func (m *Module) ExternalHandler() (handler.ExternalHandlerInterface, error) {
	if m.externalHandler == nil {
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
