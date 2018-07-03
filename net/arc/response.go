package arc

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
}
