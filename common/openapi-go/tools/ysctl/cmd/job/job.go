package job

import (
	"github.com/spf13/cobra"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

var jobLong = templates.LongDesc(`
    Job managements commands.
    
    This commands allow you to manage your jobs and list applications.`)

func NewCmdJob(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "job SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage jobs and applications",
		Long:                  jobLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.Out),
	}

	cmd.AddCommand(NewAppList(f, ioStreams))
	cmd.AddCommand(NewZoneList(f, ioStreams))
	cmd.AddCommand(NewJobList(f, ioStreams))
	cmd.AddCommand(NewJobGet(f, ioStreams))
	cmd.AddCommand(NewJobSubmit(f, ioStreams))
	cmd.AddCommand(NewFastDownload(f, ioStreams))
	return cmd
}
