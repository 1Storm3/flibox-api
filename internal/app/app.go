package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/fiber-swagger"

	_ "kinopoisk-api/docs"
	"kinopoisk-api/internal/delivery/rest"
	"kinopoisk-api/internal/metrics"
	"kinopoisk-api/internal/metrics/interceptor"
	"kinopoisk-api/shared/closer"
	"kinopoisk-api/shared/logger"
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

	a.httpServer.Get("/swagger/*", fiberSwagger.WrapHandler)

	a.httpServer.Use(interceptor.MetricsInterceptor())

	filmModule, err := a.diContainer.FilmModule()
	if err != nil {
		return err
	}
	filmHandler, err := filmModule.FilmHandler()
	if err != nil {
		return err
	}
	filmSequelModule, err := a.diContainer.FilmSequelModule()
	if err != nil {
		return err
	}

	filmSequelHandler, err := filmSequelModule.FilmSequelHandler()
	if err != nil {
		return err
	}
	userModule, err := a.diContainer.UserModule()
	if err != nil {
		return err
	}

	userHandler, err := userModule.UserHandler()
	if err != nil {
		return err
	}

	userFilmModule, err := a.diContainer.UserFilmModule()
	if err != nil {
		return err
	}

	userFilmHandler, err := userFilmModule.UserFilmHandler()
	if err != nil {
		return err
	}

	filmSimilarModule, err := a.diContainer.FilmSimilarModule()
	if err != nil {
		return err
	}

	filmSimilarHandler, err := filmSimilarModule.FilmSimilarHandler()
	if err != nil {
		return err
	}

	router := rest.NewRouter(
		filmHandler,
		filmSequelHandler,
		userHandler,
		filmSimilarHandler,
		userFilmHandler,
	)
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
