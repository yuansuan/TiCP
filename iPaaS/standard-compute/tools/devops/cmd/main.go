package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xsignal"
)

func Fatal(err interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
	os.Exit(1)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			Fatal(err)
		}
	}()

	cmdline, err := InitCmdline()
	if err != nil {
		Fatal(err)
	}

	ctx, cancel := xsignal.With(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err = cmdline.ExecuteContext(ctx); err != nil {
		Fatal(err)
	}
}
