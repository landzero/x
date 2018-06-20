package ivy

import (
	"testing"
)

func TestNewRegisterRequest(t *testing.T) {
	req, err := NewRegisterRequest("http://*.farm.landzero.net/hello/*")
	if err != nil {
		t.Error(err)
		return
	}
	if req.Method != MethodRegister {
		t.Errorf("invalid method: %s", req.Method)
	}
	if req.Host != "*.farm.landzero.net" {
		t.Errorf("invalid host: %s", req.Host)
	}
	if req.URL.Path != "/hello/*" {
		t.Errorf("invalid path: %s", req.URL.Path)
	}
}
