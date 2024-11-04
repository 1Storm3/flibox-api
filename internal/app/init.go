package app

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swaggo/fiber-swagger"

	"kbox-api/internal/delivery/rest"
	"kbox-api/internal/metrics"
	"kbox-api/internal/metrics/interceptor"
	"kbox-api/shared/httperror"
	"kbox-api/shared/logger"
)

func (a *App) initModulesAndHandlers() error {
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

	authModule, err := a.diContainer.AuthModule()
	if err != nil {
		return err
	}
	authHandler, err := authModule.AuthHandler()
	if err != nil {
		return err
	}

	externalModule, err := a.diContainer.ExternalModule()
	if err != nil {
		return err
	}
	externalHandler, err := externalModule.ExternalHandler()
	if err != nil {
		return err
	}

	router := rest.NewRouter(
		filmHandler,
		filmSequelHandler,
		userHandler,
		filmSimilarHandler,
		userFilmHandler,
		authHandler,
		externalHandler,
	)
	router.LoadRoutes(a.httpServer)

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

func (a *App) initJWT() {
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_SECRET_KEY не установлен в окружении")
	}

	a.httpServer.Use(func(c *fiber.Ctx) error {
		c.Locals("jwtKey", jwtKey)
		return c.Next()
	})
}

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = newDIContainer()

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.diContainer.Config().Env)
	return nil
}
