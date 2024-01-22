package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/api"
	"github.com/ejagombar/SpannerBackend/internal/storage"
	"github.com/ejagombar/SpannerBackend/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	cleanup, err := run(env)

	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	go func() {
		fmt.Println(app.Listen("0.0.0.0:" + env.PORT))
	}()

	return func() {
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, error) {
	db, err := storage.LoadBbolt("data", 1*time.Second)
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	spannerStore := api.NewSpannerStorage(db)
	spannerController := api.NewSpannerController(spannerStore, &env)
	// spannerStore.SaveToken(env.ACCESS_TOKEN, env.REFRESH_TOKEN, env.TOKEN_TIMEOUT)
	api.AddTodoRoutes(app, spannerController)

	return app, nil
}
