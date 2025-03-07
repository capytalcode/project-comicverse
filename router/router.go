package router

import (
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type router struct {
	cache       bool

	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (http.Handler, error) {
	if cfg.Assertions == nil {
		return nil, errors.New("assertions is nil")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	r := &router{
		cache:  !cfg.DisableCache,
		assert: cfg.Assertions,
		log:    cfg.Logger,
	}

	return r.setup(), nil
}

type Config struct {
	DisableCache bool

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

func (router *router) setup() http.Handler {
	router.assert.NotNil(router.log)

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

	r.HandleFunc("/dashboard", router.dashboard)

	return r
}

func (router *router) dashboard(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	w.WriteHeader(http.StatusOK)
}
