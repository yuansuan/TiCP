package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdVersion version
func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  `show version information`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version 1.0")
		},
	}
	return cmd
}
