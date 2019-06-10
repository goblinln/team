package web

import (
	"net/http"
)

// New returns an new web.Router object to register request handlers.
func New() *Router {
	return &Router{
		prefix:      "",
		middlewares: []Middleware{},
		children:    make(map[string]*Router),
		dispatchers: make(map[string][]*Dispatcher),
		notFound:    WrapFunc(http.NotFound),
	}
}

// Wrap http.Handler to web.Handler
func Wrap(handler http.Handler) Handler {
	return func(c *Context) {
		handler.ServeHTTP(c.rsp, c.req)
	}
}

// WrapFunc wraps original http handle function to web.Handler
func WrapFunc(handle func(w http.ResponseWriter, r *http.Request)) Handler {
	return func(c *Context) {
		handle(c.rsp, c.req)
	}
}

// Start httpd service
func Start(addr string, router *Router) error {
	return http.ListenAndServe(addr, router)
}
