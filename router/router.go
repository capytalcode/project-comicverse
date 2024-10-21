package router

import (
	"net/http"
	"path"

	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
	handlers    map[string]http.Handler
}

func NewRouter(mux ...*http.ServeMux) *Router {
	return &Router{
		http.NewServeMux(),
		[]middleware.Middleware{},
		map[string]http.Handler{},
	}
}

func (r *Router) Handle(p string, h http.Handler) {
	if sr, ok := h.(*Router); ok {
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

func (r *Router) HandleFunc(p string, hf http.HandlerFunc) {
	r.handle(p, hf)
}

func (r Router) handle(p string, h http.Handler) {
	if len(r.middlewares) > 0 {
		h = r.wrapMiddlewares(r.middlewares, h)
	}
	r.handlers[p] = h
	r.mux.Handle(p, h)
}

func (r *Router) Use(m middleware.Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r Router) wrapMiddlewares(ms []middleware.Middleware, h http.Handler) http.Handler {
	hf := h
	for _, m := range ms {
		hf = m(hf)
	}
	return hf
}
