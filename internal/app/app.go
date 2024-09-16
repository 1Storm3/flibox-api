package app

import (
	"context"
	"net"

	"github.com/gofiber/fiber/v2"

	"kinopoisk-api/internal/closer"
	"kinopoisk-api/internal/delivery/rest"
	"kinopoisk-api/internal/logger"
)

type mockFilmService struct{}

func (s *mockFilmService) One(ctx context.Context, id string) (interface{}, error) {
	return struct {
		Lolo string `json:"lolo"`
	}{
		Lolo: "1234rtt44",
	}, nil
}

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

	//handlers
	filmHandler := rest.NewFilmHandler(&mockFilmService{}) // a.diContainer.FilmService()
	// etc

	//router
	router := rest.NewRouter(filmHandler)
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
