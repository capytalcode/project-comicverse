package router

import (
	"fmt"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
)

type Router struct {
	mux         *http.ServeMux
	serveHttp   http.HandlerFunc
	middlewares []middleware.Middleware
}

type RouterOpts struct {
	Mux *http.ServeMux
}

type Route struct {
	Pattern string
	Handler http.Handler
}

func NewRouter(opts ...RouterOpts) *Router {
	var mux *http.ServeMux
	if len(opts) > 0 {
		mux = opts[0].Mux
	} else {
		mux = http.NewServeMux()
	}

	return &Router{
		mux:         mux,
		serveHttp:   mux.ServeHTTP,
		middlewares: []middleware.Middleware{},
	}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) HandleFunc(pattern string, hf http.HandlerFunc) {
	r.mux.HandleFunc(pattern, hf)
}

func (rt *Router) HandleRoutes(routes []Route) {
	for _, r := range routes {
		rt.Handle(r.Pattern, r.Handler)
	}
}

func (r *Router) AddMiddleware(m middleware.Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(rt.middlewares) > 0 {
		rt.serveHttp = rt.wrapMiddlewares(rt.middlewares, rt.serveHttp)
	}
	rt.serveHttp(w, r)
}

func (r *Router) wrapMiddlewares(ms []middleware.Middleware, h http.Handler) http.HandlerFunc {
	wh := h
	for _, m := range ms {
		wh = m(wh)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		mw := middleware.NewMiddlewaredResponse(w)
		wh.ServeHTTP(mw, r)
		if _, err := mw.ReallyWriteHeader(); err != nil {
			_, _ = w.Write(
				[]byte(
					fmt.Sprintf(
						"Error while trying to write middlewared response.\n%s",
						err.Error(),
					),
				),
			)
		}
	}
}
