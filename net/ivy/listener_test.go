package ivy

import (
	"bufio"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const listenTestPath = "/tmp/net.lanzero.ivy.test.listen"

func TestListener(t *testing.T) {
	var count uint64
	cond := sync.NewCond(&sync.Mutex{})
	// remove existed unix domain socket
	os.Remove(listenTestPath)
	// listen
	l, err := net.Listen("unix", listenTestPath)
	if err != nil {
		t.Fatal(err)
	}
	// build and run hub server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		// hijack
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
		// create and send a dummy request
		nreq, err := http.NewRequest(http.MethodGet, "http://what.farm.landzero.net/test/a", nil)
		if err != nil {
			t.Fatal(err)
		}
		if err = nreq.Write(brw); err != nil {
			t.Fatal(err)
		}
		brw.Flush()
		// read the response
		nres, err := http.ReadResponse(bufio.NewReader(brw), nreq)
		if err != nil {
			t.Fatal(err)
		}
		defer nres.Body.Close()
		buf, err := ioutil.ReadAll(nres.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(buf) != "OK" {
			t.Fatal("invalid body:", string(buf))
		}
	})
	s := http.Server{Handler: mux}
	go s.Serve(l)
	defer s.Shutdown(context.Background())

	// create the Ivy Listener
	var l2 net.Listener
	if l2, err = Listen("unix", listenTestPath, ListenConfig{
		Registration: "http://*.farm.landzero.net/test/*",
		PoolSize:     2,
	}); err != nil {
		t.Fatal(err)
	}

	// serve
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddUint64(&count, 1)
		cond.Broadcast()
		rw.Write([]byte("OK"))
	})
	s2 := http.Server{Handler: mux2}
	go s2.Serve(l2)
	defer s2.Shutdown(context.Background())

	cond.L.Lock()
	for count < 20 {
		cond.Wait()
	}
	cond.L.Unlock()
}
