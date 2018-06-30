package httpext

import (
	"io/ioutil"
	"net/http"

	"landzero.net/x/log"
)

// HealthCheck returns health check success / failed
func HealthCheck(url string) bool {
	var err error
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		log.Println("healthcheck: failed to dial:", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("healthcheck: bad status code:", resp.StatusCode)
		return false
	}
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("healthcheck: failed to read body:", err.Error())
	}
	log.Println("healthcheck: success:", string(body))
	return true
}
