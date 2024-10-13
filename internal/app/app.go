package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kinopoisk-api/internal/delivery/rest"
	"kinopoisk-api/internal/metrics"
	"kinopoisk-api/internal/metrics/interceptor"
	"kinopoisk-api/shared/closer"
	"kinopoisk-api/shared/logger"
	"log"
	"net"
	"net/http"
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

func (a *App) initHTTPServer(ctx context.Context) error {

	err := metrics.Init(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metrics: %v", err)
	}

	a.httpServer = fiber.New()

	a.httpServer.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	a.httpServer.Use(interceptor.MetricsInterceptor())

	filmService, err := a.diContainer.FilmService()
	userService, err := a.diContainer.UserService()

	filmSequelService, err := a.diContainer.FilmSequelService()
	filmSimilarService, err := a.diContainer.FilmSimilarService()

	if err != nil {
		return err
	}

	filmHandler := rest.NewFilmHandler(filmService)

	filmSequelHandler := rest.NewFilmSequelHandler(filmSequelService)

	userHandler := rest.NewUserHandler(userService)

	filmSimilarHandler := rest.NewFilmSimilarHandler(filmSimilarService)

	router := rest.NewRouter(filmHandler, filmSequelHandler, userHandler, filmSimilarHandler)
	router.LoadRoutes(a.httpServer)

	go func() {
		err = a.runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

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

func (a *App) runPrometheus() error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2020",
		Handler: mux,
	}
	log.Printf("Starting prometheus server on %s", prometheusServer.Addr)

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
