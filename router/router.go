package router

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type router struct {
	service *service.Service

	templates templates.ITemplate
	assets    fs.FS
	cache     bool

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
	if cfg.Assets == nil {
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

		templates: cfg.Templates,
		assets:    cfg.Assets,
		cache:     !cfg.DisableCache,

		assert: cfg.Assertions,
		log:    cfg.Logger,
	}

	return r.setup(), nil
}

type Config struct {
	Service *service.Service

	Templates    templates.ITemplate
	Assets       fs.FS
	DisableCache bool

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

func (router *router) setup() http.Handler {
	router.assert.NotNil(router.log)
	router.assert.NotNil(router.assets)

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

	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServerFS(router.assets)))

	r.HandleFunc("/dashboard/", router.dashboard)

	r.HandleFunc("/projects/{$}", router.projects)
	r.HandleFunc("/projects/{ID}/", router.projects)
	r.HandleFunc("/projects/{ID}/pages/{$}", router.pages)
	r.HandleFunc("/projects/{ID}/pages/{PageID}", router.pages)

	return r
}

func (router *router) dashboard(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(router.templates)
	router.assert.NotNil(router.service)
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	p, err := router.service.ListProjects()
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = router.templates.ExecuteTemplate(w, "dashboard", p)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
	}
}

func getMethod(r *http.Request) string {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return r.Method
	}

	m := r.FormValue("x-method")
	if m == "" {
		return r.Method
	}

	return strings.ToUpper(m)
}
