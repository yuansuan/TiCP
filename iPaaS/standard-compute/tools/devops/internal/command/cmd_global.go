package command

import (
	"github.com/spf13/cobra"
)

// NewRootCommand 返回根命令行
func NewRootCommand(s SingularityCommand, v VersionCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "devops",
		Short:         "A toolbox to simplify devops",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	var debug bool
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug verbosity logging")
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
	}

	cmd.AddCommand(s)
	cmd.AddCommand(v)

	return cmd
}
