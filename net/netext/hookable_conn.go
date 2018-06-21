package netext

import (
	"net"
	"sync/atomic"
	"time"
)

type hookedConn struct {
	c  net.Conn
	f  int32
	ch ConnEndHook
}

func (c *hookedConn) firstEndCalled() bool {
	return atomic.CompareAndSwapInt32(&c.f, 0, 1)
}

func (c *hookedConn) Read(b []byte) (n int, err error) {
	n, err = c.c.Read(b)
	if ch := c.ch; ch != nil && err != nil {
		if c.firstEndCalled() {
			c.ch.ConnEnded(c)
		}
	}
	return
}

func (c *hookedConn) Write(b []byte) (n int, err error) {
	n, err = c.c.Write(b)
	if ch := c.ch; ch != nil && err != nil {
		if c.firstEndCalled() {
			c.ch.ConnEnded(c)
		}
	}
	return
}

func (c *hookedConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *hookedConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *hookedConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *hookedConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c *hookedConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

func (c *hookedConn) Close() error {
	if ch := c.ch; ch != nil {
		if c.firstEndCalled() {
			c.ch.ConnEnded(c)
		}
	}
	return c.c.Close()
}

// ConnEndHook hook interface for net.Conn#Close method
type ConnEndHook interface {
	ConnEnded(c net.Conn)
}

// HookConnClose create a new net.Conn with #Close method hooked
func HookConnClose(c net.Conn, h ConnEndHook) net.Conn {
	return &hookedConn{c: c, ch: h}
}
