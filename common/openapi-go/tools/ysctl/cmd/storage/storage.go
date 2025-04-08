package storage

import (
	"github.com/spf13/cobra"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

var storageLong = templates.LongDesc(`
    Storage managements commands.

    This commands allow you to manage your storage.
	Storage_Endpoint must be set in config file.`)

func NewCmdStorage(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "storage SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage storage",
		Long:                  storageLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.Out),
	}

	cmd.AddCommand(NewStorageList(f, ioStreams))
	cmd.AddCommand(NewStorageDownload(f, ioStreams))
	cmd.AddCommand(NewStorageUpload(f, ioStreams))
	cmd.AddCommand(NewStorageQuota(f, ioStreams))
	cmd.AddCommand(NewStorageOperationLog(f, ioStreams))
	return cmd
}

func NewLiteCmdStorage(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "storage SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage storage",
		Long:                  storageLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.Out),
	}

	cmd.AddCommand(NewStorageDownload(f, ioStreams))
	cmd.AddCommand(NewStorageUpload(f, ioStreams))
	return cmd
}
