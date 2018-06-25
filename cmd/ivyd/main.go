package main

import (
	"container/list"
	"context"
	"errors"
	"flag"
	"io"
	"landzero.net/x/log"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"landzero.net/x/net/netext"

	"landzero.net/x/net/ivy"
	"landzero.net/x/os/osext"
)

var errIvyConnectionNotFound = errors.New("ivy connection not found")

var httpAddr string
var ivyAddr string

type registry struct {
	conns *list.List
	r     *sync.Mutex
}

func (r *registry) ConnEnded(c net.Conn) {
	r.r.Lock()
	defer r.r.Unlock()
	for e := r.conns.Front(); e != nil; e = e.Next() {
		if e.Value == c {
			r.conns.Remove(e)
			break
		}
	}
}

func (r *registry) add(c net.Conn, host, path string) {
	r.r.Lock()
	defer r.r.Unlock()
	r.conns.PushBack(netext.HookConnClose(c, r))
	log.Println("reg: connection added")
}

func (r *registry) take(host, path string) (c net.Conn, err error) {
	r.r.Lock()
	defer r.r.Unlock()
	e := r.conns.Front()
	if e == nil {
		err = errIvyConnectionNotFound
		return
	}
	r.conns.Remove(e)
	c = e.Value.(net.Conn)
	log.Println("reg: connection taken")
	return
}

type httpHandler struct {
	reg *registry
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// check hijackable
	hij, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to hijack connection"))
		return
	}
	// take connection and write request
	var ic net.Conn
	var err error
	if ic, err = h.reg.take(req.Host, req.URL.Path); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	ic.SetDeadline(time.Time{})
	defer ic.Close()
	if err = req.Write(ic); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	// hijack connection
	c, brw, err := hij.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to hijack connection"))
		return
	}
	c.SetDeadline(time.Time{})
	defer c.Close()
	// stream
	io.Copy(brw, ic)
	brw.Flush()
}

type ivyHandler struct {
	reg *registry
}

func (h *ivyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != ivy.MethodRegister {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid http method"))
		return
	}
	hij, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to hijack connection"))
		return
	}
	c, _, err := hij.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to hijack connection"))
		return
	}
	c.SetDeadline(time.Time{})
	h.reg.add(c, req.Host, req.URL.Path)
}

func main() {
	reg := &registry{conns: list.New(), r: &sync.Mutex{}}

	flag.StringVar(&httpAddr, "http.addr", "0.0.0.0:8080", "listening address for http")
	flag.StringVar(&ivyAddr, "ivy.addr", "127.0.0.1:8090", "listening address for ivy")
	flag.Parse()

	hs := &http.Server{Handler: &httpHandler{reg}, Addr: httpAddr}
	is := &http.Server{Handler: &ivyHandler{reg}, Addr: ivyAddr}

	go hs.ListenAndServe()
	go is.ListenAndServe()

	osext.WaitSignals(syscall.SIGINT, syscall.SIGTERM)

	hs.Shutdown(context.Background())
	is.Shutdown(context.Background())
}
