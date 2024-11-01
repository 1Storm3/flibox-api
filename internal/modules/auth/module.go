package auth

import (
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/modules/auth/handler"
	"kinopoisk-api/internal/modules/auth/service"
	"kinopoisk-api/internal/modules/user"
)

type Module struct {
	config      *config.Config
	authHandler *handler.AuthHandler
	authService handler.AuthService
	userModule  *user.Module
}

func NewAuthModule(userModule *user.Module, config *config.Config) *Module {
	return &Module{
		config:     config,
		userModule: userModule,
	}
}

func (a *Module) AuthService() (handler.AuthService, error) {
	if a.authService == nil {
		userService, err := a.userModule.UserService()
		if err != nil {
			return nil, err
		}
		a.authService = service.NewAuthService(userService)
	}
	return a.authService, nil
}

func (a *Module) AuthHandler() (*handler.AuthHandler, error) {
	if a.authHandler == nil {
		authService, err := a.AuthService()
		if err != nil {
			return nil, err
		}
		a.authHandler = handler.NewAuthHandler(authService)
	}
	return a.authHandler, nil
}
