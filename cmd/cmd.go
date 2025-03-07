package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/loreddev/x/tinyssert"
)

var (
	hostname     = flag.String("hostname", "localhost", "Host to listen to")
	port         = flag.Uint("port", 8080, "Port to be used for the server.")
	templatesDir = flag.String("templates", "", "Templates directory to be used instead of built-in ones.")
	verbose      = flag.Bool("verbose", false, "Print debug information on logs")
	dev          = flag.Bool("dev", false, "Run the server in debug mode.")
)

func init() {
	flag.Parse()
}

func Execute() {
	ctx := context.Background()

	assertions := tinyssert.NewDisabledAssertions()
	if *dev {
		assertions = tinyssert.NewAssertions()
	}

	level := slog.LevelError
	if *dev {
		level = slog.LevelDebug
	} else if *verbose {
		level = slog.LevelInfo
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	app := router.New(assertions, log, *dev)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *hostname, *port),
		Handler: app,
	}

	c, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Starting application",
			slog.String("host", *hostname),
			slog.Uint64("port", uint64(*port)),
			slog.Bool("verbose", *verbose),
			slog.Bool("development", *dev))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start application", slog.String("error", err.Error()))
		}
	}()

	<-c.Done()

	log.Info("Stopping application gracefully")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to stop application gracefully", slog.String("error", err.Error()))
	}

	log.Info("FINAL")
	os.Exit(0)
}
