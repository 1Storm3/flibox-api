package app

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/shared/closer"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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
		logger.Fatal("Ошибка при инициализации метрик: %v", zap.Error(err))
	}

	a.initFiberServer()
	a.initCORS()

	if err := a.initModulesAndHandlers(); err != nil {
		return err
	}

	go func() {
		err := a.runPrometheus()
		if err != nil {
			logger.Fatal("Ошибка при запуске метрик", zap.Error(err))
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
	logger.Info("Сервер метрик запущен на", zap.String("Адрес", prometheusServer.Addr))

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
