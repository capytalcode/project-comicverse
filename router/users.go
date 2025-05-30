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

func (c userController) login(w http.ResponseWriter, r *http.Request) {
	c.assert.NotNil(c.templates)
	c.assert.NotNil(c.service)

	if r.Method == http.MethodGet {
		err := c.templates.ExecuteTemplate(w, "login", nil)
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
		exception.BadRequest(errors.New(`missing "user" form value`)).ServeHTTP(w, r)
		return
	}
	if passwd == "" {
		exception.BadRequest(errors.New(`missing "password" form value`)).ServeHTTP(w, r)
		return
	}

	c.service.Login(user, passwd)
}
