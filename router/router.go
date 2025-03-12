package router

import (
	"errors"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type router struct {
	service *service.Service

	templates   *template.Template
	staticFiles fs.FS
	cache       bool

	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (http.Handler, error) {
	if cfg.Service == nil {
		return nil, errors.New("service is nil")
	}
	if cfg.Templates == nil {
		return nil, errors.New("templates is nil")
	}
	if cfg.StaticFiles == nil {
		return nil, errors.New("static files is nil")
	}
	if cfg.Assertions == nil {
		return nil, errors.New("assertions is nil")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	r := &router{
		service: cfg.Service,

		templates:   cfg.Templates,
		staticFiles: cfg.StaticFiles,
		cache:       !cfg.DisableCache,

		assert: cfg.Assertions,
		log:    cfg.Logger,
	}

	return r.setup(), nil
}

type Config struct {
	Service *service.Service

	Templates    *template.Template
	StaticFiles  fs.FS
	DisableCache bool

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

func (router *router) setup() http.Handler {
	router.assert.NotNil(router.log)
	router.assert.NotNil(router.staticFiles)

	log := router.log

	log.Debug("Initializing router")

	r := smalltrip.NewRouter(
		smalltrip.WithAssertions(router.assert),
		smalltrip.WithLogger(log.WithGroup("smalltrip")),
	)

	r.Use(middleware.Logger(log.WithGroup("requests")))
	if router.cache {
		r.Use(middleware.Cache())
	} else {
		r.Use(middleware.DisableCache())
	}

	r.Use(exception.PanicMiddleware())
	r.Use(exception.Middleware())

	r.Handle("/static", http.StripPrefix("/static/", http.FileServerFS(router.staticFiles)))

	r.HandleFunc("/dashboard", router.dashboard)

	r.HandleFunc("/projects/{id...}", router.projects)

	return r
}

func (router *router) dashboard(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(router.templates)
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	w.WriteHeader(http.StatusOK)
	err := router.templates.ExecuteTemplate(w, "dashboard", nil)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
	}
}
