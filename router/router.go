package router

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type router struct {
	userService  *service.User
	tokenService *service.Token

	templates templates.ITemplate
	assets    fs.FS
	cache     bool

	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (http.Handler, error) {
	if cfg.UserService == nil {
		return nil, errors.New("user service is nil")
	}
	if cfg.TokenService == nil {
		return nil, errors.New("token service is nil")
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
		userService:  cfg.UserService,
		tokenService: cfg.TokenService,

		templates: cfg.Templates,
		assets:    cfg.Assets,
		cache:     !cfg.DisableCache,

		assert: cfg.Assertions,
		log:    cfg.Logger,
	}

	return r.setup(), nil
}

type Config struct {
	UserService  *service.User
	TokenService *service.Token

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

	userController := newUserController(userControllerCfg{
		UserService:  router.userService,
		TokenService: router.tokenService,
		LoginPath:    "/login/",
		RedirectPath: "/",
		Templates:    router.templates,
		Assert:       router.assert,
	})

	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServerFS(router.assets)))

	r.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add a way to the user to bypass this check and see the landing page.
		//       Probably a query parameter to bypass like "?landing=true"
		if userController.isLogged(r) {
			err := router.templates.ExecuteTemplate(w, "dashboard", nil)
			if err != nil {
				exception.InternalServerError(err).ServeHTTP(w, r)
			}
			return
		}

		err := router.templates.ExecuteTemplate(w, "landing", nil)
		if err != nil {
			exception.InternalServerError(err).ServeHTTP(w, r)
		}
	})

	r.HandleFunc("/login/{$}", userController.login)
	r.HandleFunc("/register/{$}", userController.register)

	return r
}
