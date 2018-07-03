package arc

import "net/http"

// RouteBuilder handler builder
type RouteBuilder struct {
	Arc   *Arc
	Rules RouteRules
}

// Rule make a clone of RouteBuilder and append a RouteRule
func (h RouteBuilder) Rule(r RouteRule) RouteBuilder {
	return RouteBuilder{Arc: h.Arc, Rules: h.Rules.Add(r)}
}

// Method add method restriction
func (h RouteBuilder) Method(method ...string) RouteBuilder {
	return h.Rule(RouteMethodRule{Method: method})
}

// Get shortcut for .Method("GET").Path(...)
func (h RouteBuilder) Get(path string) RouteBuilder {
	return h.Method(http.MethodGet).Path(path)
}

// Post shortcut for .Method("POST").Path(...)
func (h RouteBuilder) Post(path string) RouteBuilder {
	return h.Method(http.MethodPost).Path(path)
}

// Put shortcut for .Method("PUT").Path(...)
func (h RouteBuilder) Put(path string) RouteBuilder {
	return h.Method(http.MethodPut).Path(path)
}

// Patch shortcut for .Method("PATCH").Path(...)
func (h RouteBuilder) Patch(path string) RouteBuilder {
	return h.Method(http.MethodPatch).Path(path)
}

// Delete shortcut for .Method("DELETE").Path(...)
func (h RouteBuilder) Delete(path string) RouteBuilder {
	return h.Method(http.MethodDelete).Path(path)
}

// Path add path restriction
func (h RouteBuilder) Path(path string) RouteBuilder {
	return h.Rule(RoutePathRule{Path: path})
}

// Host add host restriction
func (h RouteBuilder) Host(host ...string) RouteBuilder {
	return h.Rule(RouteHostRule{Host: host})
}

// Header add header restriction
func (h RouteBuilder) Header(name string, value ...string) RouteBuilder {
	return h.Rule(RouteHeaderRule{Name: name, Value: value})
}

// Use complete the handler and register to *arc.Arc, then create a empty RouteBuilder
func (h RouteBuilder) Use(handlers ...Handler) RouteBuilder {
	for _, handler0 := range handlers {
		handler := handler0
		h.Arc.Use(func(ctx *Context) {
			if h.Rules.Match(ctx.Req) {
				handler(ctx)
			} else {
				ctx.Next()
			}
		})
	}
	return h.Arc.Route()
}
