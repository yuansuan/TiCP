package cloudapp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type SoftwareGetOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SoftwareID string
}

var softwareGetExample = templates.Examples(`
    # Get software
    ysctl cloudapp softwareget --softwareID=4WoY5JA8mvE
`)

func NewSoftwareGetOptions(ioStreams clientcmd.IOStreams) *SoftwareGetOptions {
	return &SoftwareGetOptions{
		IOStreams: ioStreams,
	}
}

func NewSoftwareGet(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSoftwareGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "softwareget",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Get software",
		TraverseChildren:      true,
		Long:                  "Get software",
		Example:               softwareGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SoftwareID, "softwareID", o.SoftwareID, "Software ID")

	return cmd
}

func (o *SoftwareGetOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *SoftwareGetOptions) Run(args []string) error {
	software, err := o.hc.API.CloudApp.Software.User.Get(
		o.hc.API.CloudApp.Software.User.Get.Id(o.SoftwareID),
	)
	if err != nil {
		return err
	}
	bf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(software); err != nil {
		return err
	}
	fmt.Println(bf.String())
	return nil
}
