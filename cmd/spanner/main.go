package main

import (
	"fmt"
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/api"
	"github.com/ejagombar/SpannerBackend/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
	"os"
)

func main() {

	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
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

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		app.Listen("0.0.0.0:8080")
	}()

	// return a function to close the server and database
	return func() {
		app.Shutdown()
	}, nil

}

func buildServer(env config.EnvVars) (*fiber.App, error) {
	// init the storage

	// create the fiber app
	app := fiber.New()

	// add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	api.AddTodoRoutes(app)

	return app, nil
}
