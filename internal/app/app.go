package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/closer"
	"kinopoisk-api/internal/logger"
	"net"
)

type App struct {
	diContainer *diContainer
	*fiber.App
	route *routeConfig
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {

	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	if err := a.runApp(); err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {

	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initApp,
		a.initRoutes,
		a.initLogger,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initRoutes(_ context.Context) error {
	a.route = newRouteConfig(a.App, a.diContainer.Config())

	a.route.Setup()

	return nil
}

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = newDIContainer()

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.diContainer.Config().Env)
	return nil
}

func (a *App) initApp(ctx context.Context) error {
	a.App = fiber.New()

	_, err := a.diContainer.App(ctx)

	if err != nil {
		return err
	}
	return nil
}
func (a *App) runApp() error {

	l, err := net.Listen("tcp", a.diContainer.Config().App.HostPort())
	if err != nil {
		return err
	}

	if err := a.App.Listener(l); err != nil {
		return err
	}

	return nil
}
