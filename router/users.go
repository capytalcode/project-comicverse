package router

import (
	"errors"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type userController struct {
	assert    tinyssert.Assertions
	templates templates.ITemplate
	service   *service.UserService
}

func newUserController(
	service *service.UserService,
	templates templates.ITemplate,
	assert tinyssert.Assertions,
) userController {
	return userController{
		assert:    assert,
		templates: templates,
		service:   service,
	}
}

func (ctrl userController) login(w http.ResponseWriter, r *http.Request) {
	ctrl.assert.NotNil(ctrl.templates)
	ctrl.assert.NotNil(ctrl.service)

	if r.Method == http.MethodGet {
		err := ctrl.templates.ExecuteTemplate(w, "login", nil)
		if err != nil {
			exception.InternalServerError(err).ServeHTTP(w, r)
		}
		return
	}
	if r.Method != http.MethodPost {
		exception.MethodNotAllowed([]string{http.MethodGet, http.MethodPost}).
			ServeHTTP(w, r)
		return
	}

	user, passwd := r.FormValue("username"), r.FormValue("password")
	if user == "" {
		exception.BadRequest(errors.New(`missing "username" form value`)).ServeHTTP(w, r)
		return
	}
	if passwd == "" {
		exception.BadRequest(errors.New(`missing "password" form value`)).ServeHTTP(w, r)
		return
	}

	// TODO: Move token issuing to it's own service, make UserService.Login just return the user
	token, _, err := ctrl.service.Login(user, passwd)
	if errors.Is(err, service.ErrNotFound) {
		exception.NotFound(exception.WithError(errors.New("user not found"))).ServeHTTP(w, r)
		return
	} else if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	// TODO: harden the cookie policy to the same domain
	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Name:     "token",
		Value:    token,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ctrl userController) register(w http.ResponseWriter, r *http.Request) {
	ctrl.assert.NotNil(ctrl.templates)
	ctrl.assert.NotNil(ctrl.service)

	if r.Method == http.MethodGet {
		err := ctrl.templates.ExecuteTemplate(w, "register", nil)
		if err != nil {
			exception.InternalServerError(err).ServeHTTP(w, r)
		}
		return
	}

	if r.Method != http.MethodPost {
		exception.MethodNotAllowed([]string{http.MethodGet, http.MethodPost}).ServeHTTP(w, r)
		return
	}

	user, passwd := r.FormValue("username"), r.FormValue("password")
	if user == "" {
		exception.BadRequest(errors.New(`missing "username" form value`)).ServeHTTP(w, r)
		return
	}
	if passwd == "" {
		exception.BadRequest(errors.New(`missing "password" form value`)).ServeHTTP(w, r)
		return
	}

	_, err := ctrl.service.Register(user, passwd)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	// TODO: Move token issuing to it's own service, make UserService.Login just return the user
	token, _, err := ctrl.service.Login(user, passwd)
	if err == service.ErrNotFound {
		exception.NotFound(exception.WithError(errors.New("user not found"))).ServeHTTP(w, r)
		return
	} else if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	// TODO: harden the cookie policy to the same domain
	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Name:     "token",
		Value:    token,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ctrl userController) isLogged(r *http.Request) bool {
	// TODO: Check if token in valid (depends on token service being implemented)
	cs := r.CookiesNamed("token")
	return len(cs) > 0
}
