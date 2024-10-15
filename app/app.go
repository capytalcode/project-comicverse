package app

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"forge.capytal.company/capytalcode/project-comicverse/pages"
	devPages "forge.capytal.company/capytalcode/project-comicverse/pages/dev"
	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
)

type App struct {
	dev    bool
	port   int
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

	// if opts[0].Assets == nil {
	// d := http.Dir("./assets")
	// opts[0].Assets = d
	// }

	return &App{
		dev:    *opts[0].Dev,
		port:   *opts[0].Port,
		assets: opts[0].Assets,
	}
}

func (a *App) Run() {
	router := router.NewRouter()

	router.HandleRoutes(pages.PAGES)
	router.Handle("/assets/", a.assets)

	if a.dev {
		router.HandleRoutes(devPages.PAGES)

		router.AddMiddleware(middleware.DevMiddleware)
	}

	logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	mlogger := middleware.NewLoggerMiddleware(slog.New(logger))
	router.AddMiddleware(mlogger.Wrap)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", a.port), router); err != nil {
		log.Fatal(err)
	}
}
