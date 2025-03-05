package router

import (
	"errors"
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

func New(assertions tinyssert.Assertions, log *slog.Logger, dev bool) http.Handler {
	r := smalltrip.NewRouter(smalltrip.WithAssertions(assertions), smalltrip.WithLogger(log.WithGroup("smalltrip")))

	r.Use(middleware.Logger(log.WithGroup("requests")))
	if dev {
		log.Debug("Development mode activated, using development middleware")
		r.Use(middleware.Dev)
	} else {
		r.Use(middleware.PersistentCache())
	}
	r.Use(exception.PanicMiddleware())
	r.Use(exception.Middleware())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := templates.Templates().ExecuteTemplate(w, "test.html", nil)
		if err != nil {
			exception.InternalServerError(err).ServeHTTP(w, r)
		}
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		exception.InternalServerError(errors.New("TEST ERROR"),
			exception.WithData("test-data", "test-value"),
		).ServeHTTP(w, r)
	})
	r.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("TEST PANIC")
	})

	return r
}

func dashboard(w http.ResponseWriter, r *http.Request) {
}
