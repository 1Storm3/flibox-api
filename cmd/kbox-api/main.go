package main

import (
	"context"

	"go.uber.org/zap"

	_ "kbox-api/docs"
	"kbox-api/internal/app"
	"kbox-api/shared/logger"
)

// @title Swagger Kbox API
// @version 1.0
// @description Kbox API BFF
// @host localhost:8080
// @BasePath /api
func main() {

	ctx := context.Background()

	a, err := app.New(ctx)

	if err != nil {
		logger.Fatal("failed to init app ", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("failed to run app: ", zap.Error(err))
	}
}
