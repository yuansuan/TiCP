package cloudapp

import (
	"github.com/spf13/cobra"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

var cloudappLong = templates.LongDesc(`
    Cloudapp managements commands.

    This commands allow you to manage your cloudapps.`)

func NewCmdCloudApp(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "cloudapp SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage cloudapps",
		Long:                  cloudappLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.Out),
	}
	cmd.AddCommand(NewHardwareList(f, ioStreams))
	cmd.AddCommand(NewHardwareGet(f, ioStreams))
	cmd.AddCommand(NewSoftwareList(f, ioStreams))
	cmd.AddCommand(NewSoftwareGet(f, ioStreams))
	cmd.AddCommand(NewRemoteApp(f, ioStreams))
	cmd.AddCommand(NewSessionList(f, ioStreams))
	cmd.AddCommand(NewSessionGet(f, ioStreams))
	cmd.AddCommand(NewSessionDelete(f, ioStreams))
	cmd.AddCommand(NewSessionClose(f, ioStreams))
	cmd.AddCommand(NewSessionPowerOn(f, ioStreams))
	cmd.AddCommand(NewSessionPowerOff(f, ioStreams))
	cmd.AddCommand(NewSessionReboot(f, ioStreams))
	cmd.AddCommand(NewSessionReady(f, ioStreams))
	cmd.AddCommand(NewSessionStart(f, ioStreams))
	return cmd
}
