package osext

import (
	"os"
	"os/signal"
)

// WaitSignals wait until signals arrived
func WaitSignals(args ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, args...)
	<-c
}
