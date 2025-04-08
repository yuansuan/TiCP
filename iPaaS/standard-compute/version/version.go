package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 通过编译命令将git commit号注入二进制中
var gitCommit = "fulfilled by Makefile"

func Cmd() *cobra.Command {
	vc := &cobra.Command{
		Use:   "version",
		Short: "show standard-compute version",
		Long:  "show standard-compute version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("standard-compute version %s\n", gitCommit)
		},
	}

	return vc
}
