package main

import (
	"context"
	"diploma/internal/app"
	"log"
)

// @title Go JWT Swagger Example API
// @description This is a sample server with JWT authorization.
// @version 1.0

// @SecurityDefinitions.bearer

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
