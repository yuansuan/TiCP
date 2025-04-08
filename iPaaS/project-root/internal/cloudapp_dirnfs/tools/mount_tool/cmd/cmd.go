package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd/mount"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd/unmount"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd/version"
)

func NewCSPCtlCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "cspctl",
		Short: "cspctl  file mount and unmount operations",

		Run: runHelp,
	}

	cmds.AddCommand(version.NewCmdVersion())
	cmds.AddCommand(mount.NewCmdMount())
	cmds.AddCommand(unmount.NewCmdUnMount())

	return cmds
}

// runHelp default help command
func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
