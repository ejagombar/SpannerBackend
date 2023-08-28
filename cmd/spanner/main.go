package main

import (
	"fmt"
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/api"
	"github.com/ejagombar/SpannerBackend/internal/sessions"
	"github.com/ejagombar/SpannerBackend/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
	"os"
)

// TODO: add air, work out how to do sessions

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

	sessions.Init()

	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	return func() {
		app.Shutdown()
	}, nil

}

func buildServer(env config.EnvVars) (*fiber.App, error) {

	// create the fiber app
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	api.AddTodoRoutes(app, env)

	return app, nil
}
