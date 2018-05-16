// +build go1.3

// Copyright 2014 The Web Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package web is a high productive and modular web framework in Go.
package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"landzero.net/x/com"
	"landzero.net/x/net/web/inject"
)

var (
	// Root absolute path of work directory.
	Root string
)

const (
	// DEV environment development
	DEV = "development"
	// PROD environment production
	PROD = "production"
	// TEST environment test
	TEST = "test"
	// CridHeaderName name of correlation id header
	CridHeaderName = "X-Correlation-ID"
	// CridParamName name of correlation id parameter
	CridParamName = "_crid"
)

// Handler can be any callable function.
// Web attempts to inject services into the handler's argument list,
// and panics if an argument could not be fullfilled via dependency injection.
type Handler interface{}

// handlerFuncInvoker is an inject.FastInvoker wrapper of func(http.ResponseWriter, *http.Request).
type handlerFuncInvoker func(http.ResponseWriter, *http.Request)

func (invoke handlerFuncInvoker) Invoke(params []interface{}) ([]reflect.Value, error) {
	invoke(params[0].(http.ResponseWriter), params[1].(*http.Request))
	return nil, nil
}

// internalServerErrorInvoker is an inject.FastInvoker wrapper of func(rw http.ResponseWriter, err error).
type internalServerErrorInvoker func(rw http.ResponseWriter, err error)

func (invoke internalServerErrorInvoker) Invoke(params []interface{}) ([]reflect.Value, error) {
	invoke(params[0].(http.ResponseWriter), params[1].(error))
	return nil, nil
}

// validateAndWrapHandler makes sure a handler is a callable function, it panics if not.
// When the handler is also potential to be any built-in inject.FastInvoker,
// it wraps the handler automatically to have some performance gain.
func validateAndWrapHandler(h Handler) Handler {
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Web handler must be a callable function")
	}

	if !inject.IsFastInvoker(h) {
		switch v := h.(type) {
		case func(*Context):
			return ContextInvoker(v)
		case func(*Context, *log.Logger):
			return LoggerInvoker(v)
		case func(http.ResponseWriter, *http.Request):
			return handlerFuncInvoker(v)
		case func(http.ResponseWriter, error):
			return internalServerErrorInvoker(v)
		}
	}
	return h
}

// validateAndWrapHandlers preforms validation and wrapping for each input handler.
// It accepts an optional wrapper function to perform custom wrapping on handlers.
func validateAndWrapHandlers(handlers []Handler, wrappers ...func(Handler) Handler) []Handler {
	var wrapper func(Handler) Handler
	if len(wrappers) > 0 {
		wrapper = wrappers[0]
	}

	wrappedHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		h = validateAndWrapHandler(h)
		if wrapper != nil && !inject.IsFastInvoker(h) {
			h = wrapper(h)
		}
		wrappedHandlers[i] = h
	}

	return wrappedHandlers
}

// extractCrid extract X-Correlation-ID
func extractCrid(req *http.Request) (crid string) {
	crid = strings.TrimSpace(req.Header.Get(CridHeaderName))
	if len(crid) == 0 && req.URL != nil && req.URL.Query() != nil {
		crid = strings.TrimSpace(req.URL.Query().Get(CridParamName))
	}
	if len(crid) == 0 {
		crid = "-"
	}
	return
}

// Web represents the top level web application.
// inject.Injector methods can be invoked to map services on a global level.
type Web struct {
	// Env is the environment that Web is executing in.
	// The MACARON_ENV is read on initialization to set this variable.
	env string

	inject.Injector
	befores  []BeforeHandler
	handlers []Handler
	action   Handler

	hasURLPrefix bool
	urlPrefix    string // For suburl support.
	*Router

	logger *log.Logger
}

// NewWithLogger creates a bare bones Web instance.
// Use this method if you want to have full control over the middleware that is used.
// You can specify logger output writer with this function.
func NewWithLogger(out io.Writer) *Web {
	m := &Web{
		env:      DEV,
		Injector: inject.New(),
		action:   func() {},
		Router:   NewRouter(),
		logger:   log.New(out, log.Prefix(), log.Flags()),
	}
	m.SetEnv(os.Getenv("MACARON_ENV"))
	m.Router.m = m
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	m.NotFound(http.NotFound)
	m.InternalServerError(func(rw http.ResponseWriter, err error) {
		http.Error(rw, err.Error(), 500)
	})
	return m
}

// New creates a bare bones Web instance.
// Use this method if you want to have full control over the middleware that is used.
func New() *Web {
	return NewWithLogger(os.Stdout)
}

// Classic creates a classic Web with some basic default middleware:
// web.Logger, web.Recovery and web.Static.
func Classic() *Web {
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	return m
}

// Handlers sets the entire middleware stack with the given Handlers.
// This will clear any current middleware handlers,
// and panics if any of the handlers is not a callable function
func (m *Web) Handlers(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		m.Use(handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked.
// This is set to web.Router in a web.Classic().
func (m *Web) Action(handler Handler) {
	handler = validateAndWrapHandler(handler)
	m.action = handler
}

// BeforeHandler represents a handler executes at beginning of every request.
// Web stops future process when it returns true.
type BeforeHandler func(rw http.ResponseWriter, req *http.Request) bool

func (m *Web) Before(handler BeforeHandler) {
	m.befores = append(m.befores, handler)
}

// Use adds a middleware Handler to the stack,
// and panics if the handler is not a callable func.
// Middleware Handlers are invoked in the order that they are added.
func (m *Web) Use(handler Handler) {
	handler = validateAndWrapHandler(handler)
	m.handlers = append(m.handlers, handler)
}

func (m *Web) createContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		env:      m.env,
		Injector: inject.New(),
		handlers: m.handlers,
		action:   m.action,
		index:    0,
		Router:   m.Router,
		Req:      Request{req},
		Resp:     NewResponseWriter(rw),
		Render:   &DummyRender{rw},
		Data:     make(map[string]interface{}),
		crid:     extractCrid(req),
		logger:   m.logger,
	}
	c.SetParent(m)
	c.Map(c)
	c.MapTo(c.Resp, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// ServeHTTP is the HTTP Entry point for a Web instance.
// Useful if you want to control your own HTTP server.
// Be aware that none of middleware will run without registering any router.
func (m *Web) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.hasURLPrefix {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, m.urlPrefix)
	}
	for _, h := range m.befores {
		if h(rw, req) {
			return
		}
	}
	m.Router.ServeHTTP(rw, req)
}

func GetDefaultListenInfo() (string, int) {
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	port := com.StrTo(os.Getenv("PORT")).MustInt()
	if port == 0 {
		port = 4000
	}
	return host, port
}

// Run the http server. Listening on os.GetEnv("PORT") or 4000 by default.
func (m *Web) Run(args ...interface{}) {
	host, port := GetDefaultListenInfo()
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			host = arg
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	addr := host + ":" + com.ToStr(port)
	logger := m.GetVal(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, m.Env())
	logger.Fatalln(http.ListenAndServe(addr, m))
}

// SetURLPrefix sets URL prefix of router layer, so that it support suburl.
func (m *Web) SetURLPrefix(prefix string) {
	m.urlPrefix = prefix
	m.hasURLPrefix = len(m.urlPrefix) > 0
}

func (m *Web) SetEnv(e string) {
	if e == DEV || e == PROD || e == TEST {
		m.env = e
	} else {
		m.env = DEV
	}
}

func (m *Web) Env() string {
	return m.env
}

func init() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic("error getting work directory: " + err.Error())
	}
}
