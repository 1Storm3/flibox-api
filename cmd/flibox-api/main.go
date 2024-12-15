package main

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/app"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"go.uber.org/zap"
)

// @title Swagger Flibox API
// @version 1.0
// @description Flibox API
// @host localhost:8080
// @BasePath /api
func main() {

	ctx := context.Background()

	a, err := app.New(ctx)

	if err != nil {
		logger.Fatal("Ошибка при инициализации приложения", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("Ошибка при запуске приложения ", zap.Error(err))
	}
}
