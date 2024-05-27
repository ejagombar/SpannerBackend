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

// Load the configuration, create a SpannerStorage struct, and then run the server.
// This function also handls graceful shutdown and error handling.
func main() {
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

// Initialise and start the server, returning a cleanup function for graceful shutdown.
// The server is deployed on a go routine to allow the main process to listen for OS calls
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

// Sets up the server by loading the database to retrieve api keys.
// This database is wrapped in a SpannerStorage object which is then
// wrapped within a SpannerController struct along with the environment variables.
// This allows all methods attatched to this struct to access this data, without making it global
func buildServer(env config.EnvVars) (*fiber.App, error) {
	app := fiber.New()

	store := api.NewSpannerStorage(session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
		KeyLookup:      "cookie:session_id",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	api.AddTodoRoutes(app, env, store)

	return app, nil
}
