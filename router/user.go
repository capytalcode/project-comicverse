package router

import (
	"errors"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type userController struct {
	userSvc  *service.User
	tokenSvc *service.Token

	loginPath    string
	redirectPath string
	templates    templates.ITemplate

	assert tinyssert.Assertions
}

func newUserController(cfg userControllerCfg) userController {
	cfg.Assert.NotNil(cfg.UserService)
	cfg.Assert.NotNil(cfg.TokenService)
	cfg.Assert.NotZero(cfg.LoginPath)
	cfg.Assert.NotZero(cfg.RedirectPath)
	cfg.Assert.NotNil(cfg.Templates)

	return userController{
		userSvc:      cfg.UserService,
		tokenSvc:     cfg.TokenService,
		loginPath:    cfg.LoginPath,
		redirectPath: cfg.RedirectPath,
		templates:    cfg.Templates,
		assert:       cfg.Assert,
	}
}

type userControllerCfg struct {
	UserService  *service.User
	TokenService *service.Token

	LoginPath    string
	RedirectPath string
	Templates    templates.ITemplate

	Assert tinyssert.Assertions
}

func (ctrl userController) login(w http.ResponseWriter, r *http.Request) {
	ctrl.assert.NotNil(ctrl.templates) // TODO?: Remove these types of assertions, since golang will panic anyway
	ctrl.assert.NotNil(ctrl.userSvc)   //        when the methods of these functions are called

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

	username, passwd := r.FormValue("username"), r.FormValue("password")
	if username == "" {
		exception.BadRequest(errors.New(`missing "username" form value`)).ServeHTTP(w, r)
		return
	}
	if passwd == "" {
		exception.BadRequest(errors.New(`missing "password" form value`)).ServeHTTP(w, r)
		return
	}

	// TODO: Move token issuing to it's own service, make UserService.Login just return the user
	user, err := ctrl.userSvc.Login(username, passwd)
	if errors.Is(err, service.ErrNotFound) {
		exception.NotFound(exception.WithError(errors.New("user not found"))).ServeHTTP(w, r)
		return
	} else if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	token, err := ctrl.tokenSvc.Issue(user)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	// TODO: harden the cookie policy to the same domain
	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Name:     "authorization",
		Value:    token,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, ctrl.redirectPath, http.StatusSeeOther)
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

	username, passwd := r.FormValue("username"), r.FormValue("password")
	if username == "" {
		exception.BadRequest(errors.New(`missing "username" form value`)).ServeHTTP(w, r)
		return
	}
	if passwd == "" {
		exception.BadRequest(errors.New(`missing "password" form value`)).ServeHTTP(w, r)
		return
	}

	user, err := ctrl.userSvc.Register(username, passwd)
	if errors.Is(err, service.ErrUsernameAlreadyExists) || errors.Is(err, service.ErrPasswordTooLong) {
		exception.BadRequest(err).ServeHTTP(w, r)
		return
	} else if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	token, err := ctrl.tokenSvc.Issue(user)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	// TODO: harden the cookie policy to the same domain
	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Name:     "authorization",
		Value:    token,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ctrl userController) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ctrl.isLogged(r) {
			http.Redirect(w, r, ctrl.loginPath, http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}

var _ middleware.Middleware = userController{}.authMiddleware

func (ctrl userController) isLogged(r *http.Request) bool {
	// TODO: Check if token in valid (depends on token service being implemented)
	cs := r.CookiesNamed("token")
	return len(cs) > 0
}
