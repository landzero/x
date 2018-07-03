package arc

import (
	"landzero.net/x/container/sack"
	"net/http"
	"net/url"
)

// Handler arc handler, initialized per request
type Handler func(*Context)

// Installer arc installer, once initialized
type Installer func(*Arc)

// Arc application
type Arc struct {
	sack.Sack
	handlers []Handler // context handlers
}

// New create a new *arc.Arc
func New() (a *Arc) {
	a = &Arc{
		Sack:     sack.Sack{},
		handlers: []Handler{},
	}
	return
}

// Install install static values
func (a *Arc) Install(i Installer) {
	i(a)
}

// Use add multiple handlers
func (a *Arc) Use(handlers ...Handler) {
	a.handlers = append(a.handlers, handlers...)
}

// CreateContext create a Context from res/req
func (a *Arc) CreateContext(res http.ResponseWriter, req *http.Request) (ctx *Context) {
	ctx = &Context{
		Sack: sack.Sack{},
		Arc:  a,
		Req: &Request{
			Request:    req,
			PathParams: url.Values{},
		},
		Res: &responseWriter{ResponseWriter: res},
	}
	// inherit Sack
	ctx.SetParent(a.Sack)
	return
}

// ServeHTTP implements http.Handler
func (a *Arc) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	a.CreateContext(res, req).Next()
}

// Route create a route builder
func (a *Arc) Route() RouteBuilder {
	return RouteBuilder{Arc: a}
}
