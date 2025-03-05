package router

import (
	"errors"
	"log/slog"
	"net/http"

	"forge.capytal.company/loreddev/x/smalltrip"
	"forge.capytal.company/loreddev/x/smalltrip/exceptions"
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
	r.Use(exceptions.Middleware())

	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello world"))
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		exceptions.InternalServerError(errors.New("TEST ERROR"),
			exceptions.WithData("test-data", "test-value"),
		).ServeHTTP(w, r)
	})

	return r
}
