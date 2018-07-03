package arc

import (
	"net/http"
	"net/url"
)

// Request wraps http.Request
type Request struct {
	*http.Request

	// PathParams are ":name" noted params in path
	PathParams url.Values
}
