package main

import (
	"fmt"
	"os"
	"time"
    "crypto/tls"
	"golang.org/x/crypto/acme/autocert"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/api"
	"github.com/ejagombar/SpannerBackend/internal/storage"
	"github.com/ejagombar/SpannerBackend/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

    go func() {
        fmt.Println("Starting HTTPS server on port 443")
        ln, err := tls.Listen("tcp", ":443", getTLSConfig())
        if err != nil {
            panic(fmt.Sprintf("Error starting HTTPS server: %v", err))
        }
        if err := app.Listener(ln); err != nil {
            panic(fmt.Sprintf("Error running HTTPS server: %v", err))
        }
    }()

	return func() {
		app.Shutdown()
	}, nil
}

func getTLSConfig() *tls.Config {
    m := &autocert.Manager{
        Prompt: autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist("spanner.eagombar.uk"),
        Cache: autocert.DirCache("./certs"),
    }

    return &tls.Config{
        GetCertificate: m.GetCertificate,
        NextProtos: []string{"http/1.1", "acme-tls/1"},
    }
}

func buildServer(env config.EnvVars) (*fiber.App, error) {
	db, err := storage.LoadBbolt("bbolt_db", 1*time.Second)
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})
	spannerStore := api.NewSpannerStorage(db)
	spannerController := api.NewSpannerController(spannerStore, &env)
	api.AddTodoRoutes(app, spannerController)

	return app, nil
}
