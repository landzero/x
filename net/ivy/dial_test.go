package ivy

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

const dialTestPath = "/tmp/net.lanzero.ivy.test.dial"

func TestDial(t *testing.T) {
	os.Remove(dialTestPath)
	l, err := net.Listen("unix", dialTestPath)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Host != "*.farm.landzero.net" {
			t.Fatalf("bad host: %s", req.Host)
		}
		if req.URL.Path != "/test/*" {
			t.Fatalf("bad path: %s", req.URL.Path)
		}
		hjk, ok := rw.(http.Hijacker)
		if !ok {
			t.Fatalf("it's not hijackable")
		}
		conn, brw, err := hjk.Hijack()
		conn.SetDeadline(time.Time{})
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if _, err = brw.Write([]byte("HELLO")); err != nil {
			t.Fatal(err)
		}
		brw.Flush()
	})
	s := http.Server{Handler: mux}
	go s.Serve(l)
	defer s.Shutdown(context.Background())
	c, err := Dial("unix", dialTestPath, "http://*.farm.landzero.net/test/*")
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "HELLO" {
		t.Fatal("BAD", string(buf))
	}
}
