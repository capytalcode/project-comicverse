package router

import (
	"fmt"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	serveHttp   http.HandlerFunc
}

type RouterOpts struct {
	Mux *http.ServeMux
}

type Route struct {
	pattern string
	handler http.Handler
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
		rt.Handle(r.pattern, r.handler)
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.serveHttp(w, r)
}

