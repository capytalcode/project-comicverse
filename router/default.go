package router

import (
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
)

var DefaultRouter = NewRouter()

func Handle(pattern string, handler http.Handler) {
	DefaultRouter.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultRouter.HandleFunc(pattern, handler)
}

func HandleRoutes(routes []Route) {
	DefaultRouter.HandleRoutes(routes)
}

func Middleware(m middleware.Middleware) {
	DefaultRouter.AddMiddleware(m)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DefaultRouter.ServeHTTP(w, r)
}
