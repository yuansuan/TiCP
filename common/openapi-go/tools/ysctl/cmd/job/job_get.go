package job

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type JobGetOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	JobID string
}

var jobGetExample = templates.Examples(`
    # Get job
    ysctl job jobget --jobID=xxx
`)

func NewJobGetOptions(ioStreams clientcmd.IOStreams) *JobGetOptions {
	return &JobGetOptions{
		IOStreams: ioStreams,
	}
}

func NewJobGet(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewJobGetOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "jobget",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Get job",
		TraverseChildren:      true,
		Long:                  "Get job",
		Example:               jobGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.JobID, "jobID", "", "jobID")
	return cmd
}

func (o *JobGetOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *JobGetOptions) Run(args []string) error {
	r, err := o.hc.API.Job.JobGet(
		o.hc.API.Job.JobGet.JobId(o.JobID),
	)

	if err != nil {
		return err
	}

	b, err2 := json.MarshalIndent(r.Data, "", "    ")
	if err2 != nil {
		fmt.Println("error:", err2)
		return err2
	}
	fmt.Fprintf(o.Out, "%s\n", b)

	return nil
}
