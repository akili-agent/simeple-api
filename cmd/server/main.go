// Package main runs the simple-api HTTP service.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/routes"
)

// version is stamped into the OpenAPI document; overridden at release time
// with -ldflags "-X main.version=<tag>".
const version = "1.0.0"

// config holds the runtime configuration. Values come from the environment so
// the same binary runs unchanged in every environment.
type config struct {
	// Port is the TCP port the HTTP server listens on (PORT).
	Port int
	// LogLevel is the minimum level emitted by the logger (LOG_LEVEL).
	LogLevel string
}

// loadConfig reads configuration from the environment, falling back to
// development-friendly defaults.
func loadConfig() config {
	// goutils has no EnvInt, so parse the string form with a clear error.
	port, err := strconv.Atoi(goutils.Env("PORT", "8080"))
	if err != nil || port <= 0 || port > 65535 {
		port = 8080
	}
	return config{
		Port:     port,
		LogLevel: goutils.Env("LOG_LEVEL", "info"),
	}
}

func main() {
	if err := run(); err != nil {
		// Fatal logs at error level and exits non-zero.
		logger.Fatal("server exited", "error", err)
	}
}

func run() error {
	cfg := loadConfig()

	log := logger.New(logger.WithLevel(logger.LogLevel(cfg.LogLevel)))
	log.Info("starting simple-api", "port", cfg.Port, "logLevel", cfg.LogLevel, "version", version)

	// Okapi with conservative server timeouts; routes are declared as data in
	// internal/routes so they double as the OpenAPI definition.
	o := okapi.New(
		okapi.WithPort(cfg.Port),
		okapi.WithReadTimeout(10),
		okapi.WithWriteTimeout(10),
		okapi.WithIdleTimeout(60),
		okapi.WithLogger(log.Logger),
	)
	o.WithOpenAPIDocs(okapi.OpenAPI{
		Title:       "simple-api",
		Summary:     "A simple, production-ready Okapi API service",
		Description: "Example JSON API scaffolded with the Okapi web framework.",
		Version:     version,
		UI:          okapi.ScalarUI,
		License:     okapi.License{Name: "MIT"},
	})

	okapi.RegisterRoutes(o, routes.Routes(o))

	// Run the listener in the background so the main goroutine can wait for a
	// termination signal and trigger a graceful drain.
	errCh := make(chan error, 1)
	go func() {
		if err := o.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		// The listener failed before any signal arrived (e.g. port in use).
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	case sig := <-stop:
		log.Info("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := o.StopWithContext(ctx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	log.Info("server stopped cleanly")
	return nil
}
