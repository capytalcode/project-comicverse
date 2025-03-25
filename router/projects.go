package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
)

func (router *router) projects(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	switch getMethod(r) {
	case http.MethodGet, http.MethodHead:
		if id := r.PathValue("ID"); id != "" {
			router.getProject(w, r)
		} else {
			router.listProjects(w, r)
		}

	case http.MethodPost:
		router.createProject(w, r)

	case http.MethodDelete:
		if id := r.PathValue("ID"); id != "" {
			router.deleteProject(w, r)
		} else {
			exception.
				BadRequest(errors.New(`missing "ID" path value`)).
				ServeHTTP(w, r)
		}

	default:
		exception.MethodNotAllowed([]string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
		}).ServeHTTP(w, r)
	}
}

func (router *router) createProject(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	if getMethod(r) != http.MethodPost {
		exception.
			MethodNotAllowed([]string{http.MethodPost}).
			ServeHTTP(w, r)
		return
	}

	p, err := router.service.CreateProject()
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	router.assert.NotZero(p.ID)

	http.Redirect(w, r, fmt.Sprintf("%s/", path.Join(r.URL.Path, p.ID)), http.StatusSeeOther)
}

func (router *router) getProject(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)
	router.assert.NotNil(router.templates)

	if getMethod(r) != http.MethodGet && getMethod(r) != http.MethodHead {
		exception.
			MethodNotAllowed([]string{http.MethodGet, http.MethodHead}).
			ServeHTTP(w, r)
		return
	}

	id := r.PathValue("ID")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "ID" must be provided`)).
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

func (router *router) listProjects(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)
	router.assert.NotNil(router.templates)

	if getMethod(r) != http.MethodGet && getMethod(r) != http.MethodHead {
		exception.
			MethodNotAllowed([]string{http.MethodGet, http.MethodHead}).
			ServeHTTP(w, r)
		return
	}

	ps, err := router.service.ListProjects()
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	b, err := json.Marshal(ps)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
}

func (router *router) deleteProject(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)
	router.assert.NotNil(router.templates)

	if getMethod(r) != http.MethodDelete {
		exception.
			MethodNotAllowed([]string{http.MethodDelete}).
			ServeHTTP(w, r)
		return
	}

	id := r.PathValue("ID")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "ID" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	err := router.service.DeleteProject(id)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	err = router.templates.ExecuteTemplate(w, "partials-status", map[string]any{
		"StatusCode":      http.StatusOK,
		"Message":         fmt.Sprintf("Project %q successfully deleted", id),
		"Redirect":        "/dashboard/",
		"RedirectMessage": "Go back to dashboard",
	})
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
}
