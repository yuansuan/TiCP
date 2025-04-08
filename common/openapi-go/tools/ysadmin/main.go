package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/common/openapi-go/tools/ysadmin/impl"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ysadmin",
		Short: "远算管理员工具",
		Long:  "远算管理员工具",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			fmt.Println()
			color.New(color.FgRed).Add(color.Bold).Add(color.Italic).Println("安全生产大于一切!")
		},
	}
	rootCmd.AddCommand(impl.Cmds...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("ExecuteError: ", err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(func() {
		impl.InitConfig()
	})
}
