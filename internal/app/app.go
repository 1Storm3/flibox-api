package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"kbox-api/internal/shared/closer"
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
	if err := a.initMetrics(ctx); err != nil {
		log.Fatalf("failed to initialize metrics: %v", err)
	}

	a.initFiberServer()
	a.initCORS()

	if err := a.initModulesAndHandlers(); err != nil {
		return err
	}

	go func() {
		err := a.runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	closer.Add(func() error {
		return a.httpServer.Shutdown()
	})
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
