package main

import (
	"context"
	"go.uber.org/zap"
	"kinopoisk-api/internal/app"
	"kinopoisk-api/shared/logger"
)

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
