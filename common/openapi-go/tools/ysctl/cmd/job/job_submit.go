package job

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
)

type JobSubmitOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	JobID             string
	JsonFile          string
	PayByAccessKeyID  string
	PayByAccessSecret string
}

var jobSubmitExample = templates.Examples(`
	# Submit job
	ysctl job jobsubmit --file=job_example.json
`)

func NewJobSubmitOptions(ioStreams clientcmd.IOStreams) *JobSubmitOptions {
	return &JobSubmitOptions{
		IOStreams: ioStreams,
	}
}

func NewJobSubmit(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewJobSubmitOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "jobsubmit",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Submit job",
		TraverseChildren:      true,
		Long:                  "Submit job",
		Example:               jobSubmitExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "file", "JSON file path")

	return cmd
}

func (o *JobSubmitOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *JobSubmitOptions) Run(args []string) error {
	data, err := util.CheckReqJsonFileAndRead(o.JsonFile)
	if err != nil {
		return err
	}

	req := new(jobcreate.Request)
	err = jsoniter.Unmarshal(data, req)
	if err != nil {
		return err
	}

	r, err := o.hc.API.Job.JobCreate(
		o.hc.API.Job.JobCreate.Name(req.Name),
		o.hc.API.Job.JobCreate.Comment(req.Comment),
		o.hc.API.Job.JobCreate.Timeout(req.Timeout),
		o.hc.API.Job.JobCreate.Zone(req.Zone),
		o.hc.API.Job.JobCreate.Params(req.Params),
		o.hc.API.Job.JobCreate.ChargeParams(req.ChargeParam),
		o.hc.API.Job.JobCreate.NoRound(req.NoRound),
		o.hc.API.Job.JobCreate.AllocType(req.AllocType),
		o.hc.API.Job.JobCreate.PayBy(req.PayBy),
		o.hc.API.Job.JobCreate.PayByParams(o.PayByAccessKeyID, o.PayByAccessSecret),
	)
	if err != nil {
		return err
	}

	b, err2 := jsoniter.MarshalIndent(r, "", "    ")
	if err2 != nil {
		fmt.Println("error:", err2)
		return err2
	}
	fmt.Fprintf(o.Out, "%s\n", b)

	return nil
}
