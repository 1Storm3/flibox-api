package auth

import (
	"kbox-api/internal/config"
	"kbox-api/internal/modules/auth/handler"
	"kbox-api/internal/modules/auth/service"
	"kbox-api/internal/modules/user"
)

type ModuleInterface interface {
	AuthService() (service.AuthServiceInterface, error)
	AuthHandler() (handler.AuthHandlerInterface, error)
}

type Module struct {
	config      *config.Config
	authHandler handler.AuthHandlerInterface
	authService service.AuthServiceInterface
	userModule  user.ModuleInterface
}

func NewAuthModule(userModule user.ModuleInterface, config *config.Config) ModuleInterface {
	return &Module{
		config:     config,
		userModule: userModule,
	}
}

func (a *Module) AuthService() (service.AuthServiceInterface, error) {
	if a.authService == nil {
		userService, err := a.userModule.UserService()
		if err != nil {
			return nil, err
		}
		a.authService = service.NewAuthService(userService, a.config)
	}
	return a.authService, nil
}

func (a *Module) AuthHandler() (handler.AuthHandlerInterface, error) {
	if a.authHandler == nil {
		authService, err := a.AuthService()
		if err != nil {
			return nil, err
		}
		a.authHandler = handler.NewAuthHandler(authService)
	}
	return a.authHandler, nil
}
