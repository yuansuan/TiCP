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

type JobListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	PageOffset, PageSize int64
	JobState, Zone       string
}

var jobListExample = templates.Examples(`
    # List all jobs
    ysctl job joblist

    # List all jobs with state
    ysctl job joblist --jobState=Pending

    # List all jobs with zone
    ysctl job joblist --zone=az-zhigu

    # List all jobs with page offset and page size
    ysctl job joblist --pageOffset=0 --pageSize=100
`)

func NewJobListOptions(ioStreams clientcmd.IOStreams) *JobListOptions {
	return &JobListOptions{
		IOStreams:  ioStreams,
		PageOffset: 0,
		PageSize:   100,
	}
}

func NewJobList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewJobListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "joblist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List jobs",
		TraverseChildren:      true,
		Long:                  "List jobs",
		Example:               jobListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().Int64Var(&o.PageOffset, "pageOffset", o.PageOffset, "Page offset")
	cmd.Flags().Int64Var(&o.PageSize, "pageSize", o.PageSize, "Page size")
	cmd.Flags().StringVar(&o.JobState, "jobState", o.JobState, "Job state")
	cmd.Flags().StringVar(&o.Zone, "zone", o.Zone, "Zone")
	return cmd
}

func (o *JobListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *JobListOptions) Run(args []string) error {
	r, err := o.hc.API.Job.JobList(
		o.hc.API.Job.JobList.PageOffset(o.PageOffset),
		o.hc.API.Job.JobList.PageSize(o.PageSize),
		o.hc.API.Job.JobList.JobState(o.JobState),
		o.hc.API.Job.JobList.Zone(o.Zone),
	)
	if err != nil {
		return err
	}

	for _, job := range r.Data.Jobs {
		b, err2 := json.MarshalIndent(job, "", "    ")
		if err2 != nil {
			fmt.Println("error:", err2)
			return err2
		}
		fmt.Fprintf(o.Out, "%s\n", b)
	}
	return nil
}
