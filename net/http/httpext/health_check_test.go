package httpext

import "testing"

func TestHealthCheck(t *testing.T) {
	if !HealthCheck("http://httpbin.org/get") {
		t.Fatal("failed to check")
	}
}
