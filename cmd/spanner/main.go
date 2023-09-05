package main

import (
	"fmt"
	"os"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/api"
	"github.com/ejagombar/SpannerBackend/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
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

	store := api.NewSpannerStorage(session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
		KeyLookup:      "cookie:session_id",
	}))

	cleanup, err := run(env, store)

	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(env config.EnvVars, store *api.SpannerController) (func(), error) {
	app, err := buildServer(env, store)
	if err != nil {
		return nil, err
	}

	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	return func() {
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars, store *api.SpannerController) (*fiber.App, error) {

	// create the fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		// AllowOrigins:     "*",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	api.AddTodoRoutes(app, env, store)

	return app, nil
}
