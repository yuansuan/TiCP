package main

import (
	"os"

	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd"
)

func main() {
	command := cmd.NewDefaultYSCtlCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
