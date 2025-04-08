package job

import (
	"time"

	"github.com/spf13/cobra"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type FastDownloadOptions struct {
	clientcmd.IOStreams
	JobID         string
	LocalDir      string
	Timeout       time.Duration
	NoNeedReg     string
	NeedReg       string
	WriteSequence bool
}

var fastDownloadExample = templates.Examples(`
	# Fast Download job result
	ysctl job fast_download --job_id=I24XKL98 --local_dir="C:/abcd/" --timeout=3600s --endpoint=openapi4.yuansuan.com
        --access_id="fadsfe" --access_secret="d23d" --no_need="*.log" --proxy=""
`)

/*
退出码
0： 下载正常
1： 参数错误
2： 超时
3： 下载出错
*/

func NewFastDownloadOptions(ioStreams clientcmd.IOStreams) *FastDownloadOptions {
	return &FastDownloadOptions{
		IOStreams: ioStreams,
	}
}

func NewFastDownload(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewFastDownloadOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "fast_download",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Fast download job result",
		TraverseChildren:      true,
		Long:                  "Fast download job result",
		Example:               fastDownloadExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.JobID, "job_id", "", "Job ID")
	cmd.Flags().StringVar(&o.LocalDir, "local_dir", "./", "dir to save result")
	cmd.Flags().StringVar(&o.NoNeedReg, "no_need", "", "no need file path with regex format")
	cmd.Flags().StringVar(&o.NeedReg, "need", "", "need file path with regex format")
	cmd.Flags().DurationVar(&o.Timeout, "timeout", 0, "timeout")
	cmd.Flags().BoolVar(&o.WriteSequence, "write_sequence", false, "true will write with sequence")
	return cmd
}

func (o *FastDownloadOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	return nil
}
