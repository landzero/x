package internal

import (
	"fmt"
	"landzero.net/x/log"
)

var Logger *log.Logger

func Logf(s string, args ...interface{}) {
	if Logger == nil {
		return
	}
	Logger.Output(2, fmt.Sprintf(s, args...))
}
