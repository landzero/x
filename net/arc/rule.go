package arc

import (
	"net/url"
	"strings"
)

// RouteRule is a single rule in a route
type RouteRule interface {
	// Match match the Request, returns if matched
	Match(req *Request) bool
}

// RouteRules a slice of RouteRule
type RouteRules []RouteRule

// Add create a copy of RouteRules with new rule added
func (rs RouteRules) Add(rule RouteRule) RouteRules {
	if rs == nil {
		return RouteRules{rule}
	}
	return RouteRules{rs, rule}
}

// Match match all rules against request
func (rs RouteRules) Match(req *Request) bool {
	for _, r := range rs {
		if !r.Match(req) {
			return false
		}
	}
	return true
}

// RouteMethodRule route rule with method restriction
type RouteMethodRule struct {
	Method []string
}

// Match implements RouteRule
func (r RouteMethodRule) Match(req *Request) bool {
	for _, m := range r.Method {
		if m == req.Method {
			return true
		}
	}
	return false
}

// RoutePathRule route rule with path restriction and params extraction
type RoutePathRule struct {
	Path string
}

// Match implements RouteRule
func (r RoutePathRule) Match(req *Request) bool {
	pp := url.Values{}
	ns := sanitizePathComponents(strings.Split(r.Path, "/"))
	hs := sanitizePathComponents(strings.Split(req.URL.Path, "/"))

	// length mismatch
	if len(ns) != len(hs) {
		if len(hs) > len(ns) && len(ns) > 0 && strings.HasPrefix(ns[len(ns)-1], "*") {
			// continue if path components longer than pattern components and pattern components has a wildcard ending
		} else {
			return false
		}
	}

	// iterate pattern components
	for i, n := range ns {
		h := hs[i]
		if strings.HasPrefix(n, ":") {
			// single capture
			pp.Set(n[1:], h)
		} else if strings.HasPrefix(n, "*") {
			// wildcard capture
			pp.Set(n[1:], strings.Join(hs[i:], "/"))
			break
		} else {
			// match path component and pattern component
			if n != h {
				return false
			}
		}
	}

	// assign PathParams
	req.PathParams = pp
	return true
}

func sanitizePathComponents(in []string) []string {
	ret := make([]string, 0, len(in))
	for _, c := range in {
		if len(c) > 0 {
			ret = append(ret, c)
		}
	}
	return ret
}

// RouteHeaderRule route rule with header restriction
type RouteHeaderRule struct {
	Name  string
	Value []string
}

// Match implements RouteRule
func (r RouteHeaderRule) Match(req *Request) bool {
	v := req.Header.Get(r.Name)
	for _, vv := range r.Value {
		if vv == v {
			return true
		}
	}
	return false
}

// RouteHostRule route rule with host restriction
type RouteHostRule struct {
	Host []string
}

// Match implements RouteRule
func (r RouteHostRule) Match(req *Request) bool {
	for _, h := range r.Host {
		if h == req.Host {
			return true
		}
	}
	return false
}
