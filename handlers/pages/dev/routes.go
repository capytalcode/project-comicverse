package dev

import (
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/lib/router"
)

func Routes() router.Router {
	r := router.NewRouter()
	r.Handle("/colors", &Colors{})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("hello world"))
	})

	return r
}
