package auth

import (
	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/modules/external"
	"github.com/1Storm3/flibox-api/internal/modules/user"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	Service() (ServiceInterface, error)
	Handler() (HandlerInterface, error)
}

type Module struct {
	config         *config.Config
	handler        HandlerInterface
	service        ServiceInterface
	userModule     user.ModuleInterface
	externalModule external.ModuleInterface
}

func NewAuthModule(userModule user.ModuleInterface,
	externalModule external.ModuleInterface,
	config *config.Config,
) *Module {
	return &Module{
		config:         config,
		userModule:     userModule,
		externalModule: externalModule,
	}
}

func (a *Module) Service() (ServiceInterface, error) {
	if a.service == nil {
		userService, err := a.userModule.Service()
		if err != nil {
			return nil, err
		}
		externalService, err := a.externalModule.EmailService()
		if err != nil {
			return nil, err
		}
		a.service = NewAuthService(userService, externalService, a.config)
	}
	return a.service, nil
}

func (a *Module) Handler() (HandlerInterface, error) {
	if a.handler == nil {
		authService, err := a.Service()
		if err != nil {
			return nil, err
		}
		a.handler = NewAuthHandler(authService)
	}
	return a.handler, nil
}
