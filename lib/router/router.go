package router

import (
	"net/http"
	"path"

	"forge.capytal.company/capytalcode/project-comicverse/lib/middleware"
)

type Router interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
	Use(middleware middleware.Middleware)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type defaultRouter struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
	handlers    map[string]http.Handler
}

func NewRouter(mux ...*http.ServeMux) Router {
	return &defaultRouter{
		http.NewServeMux(),
		[]middleware.Middleware{},
		map[string]http.Handler{},
	}
}

func (r *defaultRouter) Handle(p string, h http.Handler) {
	if sr, ok := h.(*defaultRouter); ok {
		for sp, sh := range sr.handlers {
			wh := sh
			if len(sr.middlewares) > 0 {
				wh = sr.wrapMiddlewares(sr.middlewares, wh)
			}
			r.handle(path.Join(p, sp), wh)
		}
	} else {
		r.handle(p, h)
	}
}

func (r *defaultRouter) HandleFunc(p string, hf http.HandlerFunc) {
	r.handle(p, hf)
}

func (r defaultRouter) handle(p string, h http.Handler) {
	if len(r.middlewares) > 0 {
		h = r.wrapMiddlewares(r.middlewares, h)
	}
	r.handlers[p] = h
	r.mux.Handle(p, h)
}

func (r *defaultRouter) Use(m middleware.Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r defaultRouter) wrapMiddlewares(ms []middleware.Middleware, h http.Handler) http.Handler {
	hf := h
	for _, m := range ms {
		hf = m(hf)
	}
	return hf
}
