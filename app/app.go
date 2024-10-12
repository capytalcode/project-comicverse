package app

import (
	"fmt"
	"log"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/pages"
	"forge.capytal.company/capytalcode/project-comicverse/router"
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

	if opts[0].Assets == nil {
		// d := http.Dir("./assets")
		// opts[0].Assets = d
	}

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

	if err := http.ListenAndServe(fmt.Sprintf(":%v", a.port), router); err != nil {
		log.Fatal(err)
	}
}
