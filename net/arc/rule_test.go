package arc

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRouteHeaderRule_Match(t *testing.T) {
	r := RouteHeaderRule{Name: "X-Test-Header", Value: []string{"AAA", "BBB"}}
	req := &Request{Request: &http.Request{Header: http.Header{}}}
	req.Header.Set("X-Test-Header", "AAA")
	if !r.Match(req) {
		t.Error("header rule failed 1")
	}
	req.Header.Set("X-Test-Header", "BBB")
	if !r.Match(req) {
		t.Error("header rule failed 2")
	}
	req.Header.Set("X-Test-Header", "CCC")
	if r.Match(req) {
		t.Error("header rule failed 3")
	}
}

func TestRouteHostRule_Match(t *testing.T) {
	r := RouteHostRule{Host: []string{"host1.landzero.net", "host2.landzero.net"}}
	req := &Request{Request: &http.Request{Host: "host1.landzero.net"}}
	if !r.Match(req) {
		t.Error("host rule failed 1")
	}
	req.Host = "host2.landzero.net"
	if !r.Match(req) {
		t.Error("host rule failed 2")
	}
	req.Host = "islandzero.net"
	if r.Match(req) {
		t.Error("host rule failed 3")
	}
}

func TestRouteMethodRule_Match(t *testing.T) {
	r := RouteMethodRule{Method: []string{http.MethodGet, http.MethodPost}}
	req := &Request{Request: &http.Request{}}
	req.Method = http.MethodGet
	if !r.Match(req) {
		t.Error("method rule failed 1")
	}
	req.Method = http.MethodPost
	if !r.Match(req) {
		t.Error("method rule failed 2")
	}
	req.Method = http.MethodPut
	if r.Match(req) {
		t.Error("method rule failed 3")
	}
}

func TestRoutePathRule_Match(t *testing.T) {
	r := RoutePathRule{Path: "/hello/world"}
	req := &Request{Request: &http.Request{}, PathParams: url.Values{}}
	req.URL, _ = url.Parse("/hello/world")
	if !r.Match(req) {
		t.Error("path rule failed 1")
	}
	req.URL, _ = url.Parse("/hello/world1")
	if r.Match(req) {
		t.Error("path rule failed 2")
	}
	req.URL, _ = url.Parse("/hello/world/1")
	if r.Match(req) {
		t.Error("path rule failed 3")
	}
	r = RoutePathRule{Path: "hello/world/:name"}
	if !r.Match(req) {
		t.Error("path rule failed 4")
	} else {
		if req.PathParams.Get("name") != "1" {
			t.Error("path rule failed 5")
		}
	}
	req.URL, _ = url.Parse("/hello/world/1/bug")
	if r.Match(req) {
		t.Error("path rule failed 6")
	}
	r = RoutePathRule{Path: "hello/world/*any"}
	req.URL, _ = url.Parse("/hello")
	if r.Match(req) {
		t.Error("path rule failed 7")
	}
	req.URL, _ = url.Parse("/hello/world/some/thing/wired")
	if !r.Match(req) {
		t.Error("path rule failed 8")
	} else {
		if req.PathParams.Get("any") != "some/thing/wired" {
			t.Error("path rule failed 9")
		}
	}
	r = RoutePathRule{Path: "hello/:name/*any"}
	req.URL, _ = url.Parse("/hello/world/something/wired")
	if !r.Match(req) {
		t.Error("path rule failed 10")
	} else {
		if req.PathParams.Get("any") != "something/wired" || req.PathParams.Get("name") != "world" {
			t.Error("path rule failed 11")
		}
	}
}

func TestRouteRules_Add(t *testing.T) {
	// nil should work
	var r RouteRules
	r = r.Add(RouteMethodRule{Method: []string{http.MethodGet, http.MethodPost}})
	req := &Request{Request: &http.Request{}}
	req.Method = http.MethodGet
	if !r.Match(req) {
		t.Error("add failed 1")
	}
	req.Method = http.MethodPost
	if !r.Match(req) {
		t.Error("add failed 2")
	}
	req.Method = http.MethodPut
	if r.Match(req) {
		t.Error("add failed 3")
	}
	// multiple
	r = r.Add(RouteHostRule{Host: []string{"host1.landzero.net", "host2.landzero.net"}})
	req.Method = http.MethodGet
	req.Host = "host2.landzero.net"
	if !r.Match(req) {
		t.Error("add failed 4")
	}
	req.Method = http.MethodPut
	if r.Match(req) {
		t.Error("add failed 5")
	}
	req.Method = http.MethodGet
	req.Host = "host3.landzero.net"
	if r.Match(req) {
		t.Error("add failed 6")
	}
}

func TestRouteRules_Match(t *testing.T) {
	// already tested
}
