package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"kinopoisk-api/internal/delivery/rest"
	"kinopoisk-api/shared/closer"
	"kinopoisk-api/shared/logger"
	"net"
)

type App struct {
	diContainer *diContainer
	httpServer  *fiber.App
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

	if err := a.runHTTPServer(); err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {

	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = fiber.New()

	a.httpServer.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	filmService, err := a.diContainer.FilmService()
	sequelService, err := a.diContainer.SequelService()
	if err != nil {
		return err
	}
	filmHandler := rest.NewFilmHandler(filmService)
	sequelHandler := rest.NewSequelHandler(sequelService)
	router := rest.NewRouter(filmHandler, sequelHandler)
	router.LoadRoutes(a.httpServer)
	closer.Add(func() error {
		return a.httpServer.Shutdown()
	})
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
func (a *App) runHTTPServer() error {

	l, err := net.Listen("tcp", a.diContainer.Config().App.HostPort())
	if err != nil {
		return err
	}

	if err := a.httpServer.Listener(l); err != nil {
		return err
	}

	return nil
}
