package account

import (
	"github.com/spf13/cobra"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

var accountLong = templates.LongDesc(`
    Account managements commands.
    
    This commands allow you to manage your account.`)

func NewCmdAccount(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "account COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage account",
		Long:                  accountLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.Out),
	}

	cmd.AddCommand(NewUserAdd(f, ioStreams))
	cmd.AddCommand(NewExchange(f, ioStreams))
	return cmd
}
