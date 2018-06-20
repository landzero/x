package ivy

import (
	"net/http"
)

// MethodRegister the HTTP method REGISTER
const MethodRegister = "REGISTER"

// NewRegisterRequest create a new consume request
func NewRegisterRequest(url string) (*http.Request, error) {
	return http.NewRequest(MethodRegister, url, nil)
}
