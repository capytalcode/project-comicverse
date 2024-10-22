package pages

import (
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/capytalcode/project-comicverse/router/rerrors"
)

func Routes(logger *slog.Logger) router.Router {
	r := router.NewRouter()

	mErrors := rerrors.NewErrorMiddleware(ErrorPage{}.Component, logger)
	r.Use(mErrors.Wrap)

	r.Handle("/dashboard", &Dashboard{})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			rerrors.NotFound().ServeHTTP(w, r)
		}

	})

	return r
}
