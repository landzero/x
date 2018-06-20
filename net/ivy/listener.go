package ivy

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"landzero.net/x/net/netext"
)

type listener struct {
	network string
	address string
	config  ListenConfig
	count   uint64
	cond    *sync.Cond
	closed  bool
}

// ConnCloseCalled implements netext.ConnCloseHook
func (l *listener) ConnCloseCalled(c net.Conn, first bool) {
	if first {
		// decrease count
		atomic.AddUint64(&l.count, ^uint64(0))
		// signal the cond
		l.cond.Signal()
	}
}

func (l *listener) Accept() (c net.Conn, err error) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	// wait until count exceeded or listener closed
	for l.count > l.config.PoolSize && !l.closed {
		l.cond.Wait()
	}
	// just return if closed
	if l.closed {
		return nil, ErrListenerClosed
	}
	// dial
	if c, err = Dial(l.network, l.address, l.config.Registration); err != nil {
		return
	}
	// increase count
	atomic.AddUint64(&l.count, 1)
	// hook net.Conn#Close
	c = netext.HookConnClose(c, l)
	return
}

func (l *listener) Addr() net.Addr {
	return &net.IPAddr{}
}

func (l *listener) Close() error {
	// mark listener closed
	l.closed = true
	// notify running Accept() loop
	l.cond.Broadcast()
	return nil
}

var (
	// ErrListenerClosed error listener is closed
	ErrListenerClosed = errors.New("listener closed")
)

// ListenConfig config for Listen()
type ListenConfig struct {
	// Registration registration as a url
	Registration string
	// PoolSize size of pool
	PoolSize uint64
}

// Listen register on an IvyHub and returns a virtual net.Listener
func Listen(network, address string, cfg ListenConfig) (net.Listener, error) {
	// cfg.PoolSize must be greater than 0
	if cfg.PoolSize < 1 {
		cfg.PoolSize = 5
	}
	return &listener{
		network: network,
		address: address,
		config:  cfg,
		cond:    sync.NewCond(&sync.Mutex{}),
	}, nil
}
