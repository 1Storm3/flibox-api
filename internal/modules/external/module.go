package external

import "github.com/1Storm3/flibox-api/internal/config"

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	S3Service() (S3ServiceInterface, error)
	EmailService() (EmailServiceInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	cfg          *config.Config
	service      ServiceInterface
	s3Service    S3ServiceInterface
	emailService EmailServiceInterface
	handler      HandlerInterface
}

func NewExternalModule(cfg *config.Config) *Module {
	return &Module{
		cfg: cfg,
	}
}

func (m *Module) Service() (ServiceInterface, error) {
	if m.service == nil {
		m.service = NeewExternalService(m.cfg)
	}
	return m.service, nil
}

func (m *Module) S3Service() (S3ServiceInterface, error) {
	if m.s3Service == nil {

		s3Service, err := NewS3Service(m.cfg)
		if err != nil {
			return nil, err
		}
		m.s3Service = s3Service
	}
	return m.s3Service, nil
}

func (m *Module) EmailService() (EmailServiceInterface, error) {
	if m.emailService == nil {
		m.emailService = NewEmailService(m.cfg)
	}
	return m.emailService, nil
}

func (m *Module) Handler() (HandlerInterface, error) {
	if m.handler == nil {
		externalService, err := m.Service()
		if err != nil {
			return nil, err
		}
		s3Service, err := m.S3Service()
		if err != nil {
			return nil, err
		}

		m.handler = NewExternalHandler(externalService, s3Service)
	}
	return m.handler, nil
}
