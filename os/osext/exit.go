// osext

package osext

import (
	"os"
)

// ExitCode exit code will be used when calling DoExit()
var ExitCode int

// DoExit execute a os.Exit() with previously set error code
// basically this should be defered in first line of main function
func DoExit() {
	os.Exit(ExitCode)
}

// WillExit set the error code will be used when calling DoExit()
func WillExit(exitCode int) {
	ExitCode = exitCode
}
