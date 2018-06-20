package netext

import (
	"net"
	"sync/atomic"
)

type closeHookedConn struct {
	net.Conn
	f  int32
	ch ConnCloseHook
}

func (c *closeHookedConn) Close() error {
	c.ch.ConnCloseCalled(c, atomic.CompareAndSwapInt32(&c.f, 0, 1))
	return c.Conn.Close()
}

// ConnCloseHook hook interface for net.Conn#Close method
type ConnCloseHook interface {
	ConnCloseCalled(c net.Conn, first bool)
}

// HookConnClose create a new net.Conn with #Close method hooked
func HookConnClose(c net.Conn, h ConnCloseHook) net.Conn {
	return &closeHookedConn{Conn: c, ch: h}
}
