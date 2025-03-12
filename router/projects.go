package router

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
)

func (router *router) projects(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	id := r.PathValue("id")
	if id != "" {
		router.getProject(w, r)
		return
	}

	router.createProject(w, r)
}

func (router *router) createProject(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	if r.Method != http.MethodPost {
		exception.
			MethodNotAllowed([]string{http.MethodPost}).
			ServeHTTP(w, r)
		return
	}

	router.log.Debug("Creating new project", slog.Any("servce", router.service))
	p, err := router.service.CreateProject()
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	router.log.Debug("New project created", slog.String("id", p.ID))
	router.assert.NotZero(p.ID)

	http.Redirect(w, r, path.Join(r.URL.Path, p.ID), http.StatusSeeOther)
}

func (router *router) getProject(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)
	router.assert.NotNil(router.templates)

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		exception.
			MethodNotAllowed([]string{http.MethodGet, http.MethodHead}).
			ServeHTTP(w, r)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "id" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	p, err := router.service.GetProject(id)
	switch {
	case errors.Is(err, service.ErrProjectNotExists):
		exception.NotFound().ServeHTTP(w, r)
		return

	case err != nil:
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	err = router.templates.ExecuteTemplate(w, "project", p)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
}
