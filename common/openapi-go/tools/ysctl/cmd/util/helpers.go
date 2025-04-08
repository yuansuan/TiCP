package util

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	// DefaultErrorExitCode defines the default exit code.
	DefaultErrorExitCode = 1
)

var fatalErrHandler = fatal

func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

func checkErr(err error, handleErr func(string, int)) {

	if err == nil {
		return
	}
	msg := err.Error()
	handleErr(msg, DefaultErrorExitCode)

}

func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

var ErrExit = fmt.Errorf("exit")

func DefaultSubCommandRun(out io.Writer) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		c.SetOutput(out)
		RequireNoArguments(c, args)
		c.Help()
		CheckErr(ErrExit)
	}
}

func RequireNoArguments(c *cobra.Command, args []string) {
	if len(args) > 0 {
		CheckErr(UsageErrorf(c, "unknown command %q", strings.Join(args, " ")))
	}
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func CheckReqJsonFileAndRead(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("request json file cannot be empty")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read request json file [%s] failed, %w", file, err)
	}

	return data, nil
}
