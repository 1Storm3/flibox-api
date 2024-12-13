package app

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swaggo/fiber-swagger"

	"kbox-api/internal/delivery/middleware"
	"kbox-api/internal/delivery/rest"
	"kbox-api/internal/metrics"
	"kbox-api/internal/metrics/interceptor"
	"kbox-api/internal/shared/httperror"
	"kbox-api/internal/shared/logger"
)

func (a *App) initModulesAndHandlers() error {
	filmModule, err := a.diContainer.FilmModule()
	if err != nil {
		return err
	}
	filmHandler, err := filmModule.Handler()
	if err != nil {
		return err
	}
	filmSequelModule, err := a.diContainer.FilmSequelModule()
	if err != nil {
		return err
	}
	filmSequelHandler, err := filmSequelModule.Handler()
	if err != nil {
		return err
	}
	userModule, err := a.diContainer.UserModule()
	if err != nil {
		return err
	}
	userHandler, err := userModule.Handler()
	if err != nil {
		return err
	}
	userFilmModule, err := a.diContainer.UserFilmModule()
	if err != nil {
		return err
	}
	userFilmHandler, err := userFilmModule.Handler()
	if err != nil {
		return err
	}
	filmSimilarModule, err := a.diContainer.FilmSimilarModule()
	if err != nil {
		return err
	}
	filmSimilarHandler, err := filmSimilarModule.Handler()
	if err != nil {
		return err
	}
	authModule, err := a.diContainer.AuthModule()
	if err != nil {
		return err
	}
	authHandler, err := authModule.Handler()
	if err != nil {
		return err
	}
	externalModule, err := a.diContainer.ExternalModule()
	if err != nil {
		return err
	}
	externalHandler, err := externalModule.Handler()
	if err != nil {
		return err
	}
	commentModule, err := a.diContainer.CommentModule()
	if err != nil {
		return err
	}
	commentHandler, err := commentModule.Handler()
	if err != nil {
		return err
	}
	collectionModule, err := a.diContainer.CollectionModule()
	if err != nil {
		return err
	}
	collectionHandler, err := collectionModule.Handler()
	if err != nil {
		return err
	}
	collectionFilmModule, err := a.diContainer.CollectionFilmModule()
	if err != nil {
		return err
	}
	collectionFilmHandler, err := collectionFilmModule.Handler()
	if err != nil {
		return err
	}
	historyFilmsModule, err := a.diContainer.HistoryFilmsModule()
	if err != nil {
		return err
	}
	historyFilmsHandler, err := historyFilmsModule.Handler()
	if err != nil {
		return err
	}
	userRepo, err := userModule.Repository()
	if err != nil {
		return err
	}

	config := a.diContainer.Config()

	authMiddleware := middleware.AuthMiddleware(userRepo, config)

	router := rest.NewRouter(
		filmHandler,
		filmSequelHandler,
		userHandler,
		filmSimilarHandler,
		userFilmHandler,
		authHandler,
		externalHandler,
		commentHandler,
		collectionHandler,
		collectionFilmHandler,
		historyFilmsHandler,
	)
	router.LoadRoutes(a.httpServer, authMiddleware)

	a.httpServer.Get("/swagger/*", fiberSwagger.WrapHandler)
	a.httpServer.Use(interceptor.MetricsInterceptor())

	return nil
}

func (a *App) initCORS() {
	a.httpServer.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
}

func (a *App) initMetrics(ctx context.Context) error {
	return metrics.Init(ctx)
}

func (a *App) initFiberServer() {
	a.httpServer = fiber.New(fiber.Config{
		ErrorHandler: a.customErrorHandler(),
	})
}

func (a *App) customErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var message string

		var httpErr *httperror.Error
		if errors.As(err, &httpErr) {
			code = httpErr.Code()
			message = httpErr.Error()
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}

		return ctx.Status(code).JSON(fiber.Map{
			"statusCode": code,
			"message":    message,
		})
	}
}

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = newDIContainer()

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.diContainer.Config().Env)
	return nil
}
