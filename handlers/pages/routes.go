package pages

import (
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/lib/router"
	"forge.capytal.company/capytalcode/project-comicverse/lib/router/rerrors"
)

func Routes(logger *slog.Logger) router.Router {
	r := router.NewRouter()

	r.Use(rerrors.NewErrorMiddleware(ErrorPage{}.Component, logger))

	r.Handle("/dashboard", &Dashboard{})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			rerrors.NotFound().ServeHTTP(w, r)
		}
	})

	return r
}
