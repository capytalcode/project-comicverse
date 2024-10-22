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

	"forge.capytal.company/capytalcode/project-comicverse/configs"
	"forge.capytal.company/capytalcode/project-comicverse/handlers/pages"
	devPages "forge.capytal.company/capytalcode/project-comicverse/handlers/pages/dev"
	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
)

type App struct {
	dev    bool
	port   int
	assets http.Handler
	logger *slog.Logger
	server *http.Server
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

	if opts[0].Assets == nil {
		d := http.FileServer(http.Dir("./assets"))
		opts[0].Assets = d
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
	mlogger := middleware.NewLoggerMiddleware(a.logger)

	r := router.NewRouter()

	r.Use(mlogger.Wrap)

	if configs.DEVELOPMENT {
		a.logger.Info("RUNNING IN DEVELOPMENT MODE")

		r.Use(middleware.DevMiddleware)
		r.Handle("/_dev", devPages.Routes())

	} else {
		r.Use(middleware.CacheMiddleware)
	}

	r.Handle("/assets/", a.assets)
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
