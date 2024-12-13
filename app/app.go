package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forge.capytal.company/capytalcode/project-comicverse/assets"
	"forge.capytal.company/capytalcode/project-comicverse/configs"
	"forge.capytal.company/capytalcode/project-comicverse/handlers/pages"
	devPages "forge.capytal.company/capytalcode/project-comicverse/handlers/pages/dev"
	"forge.capytal.company/capytalcode/project-comicverse/lib/middleware"
	"forge.capytal.company/capytalcode/project-comicverse/lib/router"
)

type App struct {
	dev    bool
	port   int
	logger *slog.Logger
	server *http.Server
	assets http.Handler
}

type AppOpts struct {
	Dev    *bool
	Port   *int
	Assets http.Handler
}

func NewApp(opts ...AppOpts) *App {
	if len(opts) == 0 {
		opts[0] = AppOpts{}
	}

	if opts[0].Dev == nil {
		d := false
		opts[0].Dev = &d
	}

	if opts[0].Port == nil {
		d := 8080
		opts[0].Port = &d
	}

	app := &App{
		dev:    *opts[0].Dev,
		port:   *opts[0].Port,
		assets: opts[0].Assets,
	}

	configs.DEVELOPMENT = app.dev

	app.setLogger()
	app.setServer()

	return app
}

func (a *App) setLogger() {
	a.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func (a *App) setServer() {
	r := router.NewRouter()

	r.Use(middleware.NewLoggerMiddleware(a.logger))

	if configs.DEVELOPMENT {
		a.logger.Info("RUNNING IN DEVELOPMENT MODE")

		r.Use(middleware.DevMiddleware)
		r.Handle("/_dev", devPages.Routes())
	} else {
		r.Use(middleware.CacheMiddleware)
	}

	if configs.DEVELOPMENT && a.assets != nil {
		r.Handle("/assets/", a.assets)
	} else {
		r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServerFS(assets.ASSETS)))
	}

	r.Handle("/", pages.Routes(a.logger))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", a.port),
		Handler: r,
	}

	a.server = &srv
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Listen and serve returned error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()
	a.logger.Info("Gracefully shutting doing server")
	if err := a.server.Shutdown(context.TODO()); err != nil {
		a.logger.Error("Server shut down returned an error", slog.String("error", err.Error()))
	}

	a.logger.Info("FINAL")
}
