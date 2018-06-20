package ivy

import (
	"net"
	"net/http"
)

// Dial dial a single connection to an IvyHub
func Dial(network, address, registration string) (c net.Conn, err error) {
	// create REGISTER request
	var req *http.Request
	if req, err = NewRegisterRequest(registration); err != nil {
		return
	}
	// dial
	if c, err = net.Dial(network, address); err != nil {
		return
	}
	// send REGISTER request
	if err = req.Write(c); err != nil {
		c.Close()
		c = nil
		return
	}
	return
}
