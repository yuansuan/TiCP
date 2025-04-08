package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DefaultErrorExitCode defines the default exit code.
	DefaultErrorExitCode = 1
)

// default error handle
var fatalErrHandler = fatal

// fatal error handler
func fatal(msg string, code int) {
	if len(msg) > 0 {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// CheckErr checks
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// checkErr checks error
func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}
	s := err.Error()
	handleErr(s, DefaultErrorExitCode)
}
