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

type SoftwareListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	PageOffset, PageSize int
	Name, Platform, Zone string
}

var softwareListExample = templates.Examples(`
    # List all softwares
    ysctl cloudapp softwarelist

    # List all softwares with page offset and page size
    ysctl cloudapp softwarelist --pageOffset=0 --pageSize=100

    # List all softwares with name
    ysctl cloudapp softwarelist --name=ys-1c1g1g1s

    # List all softwares with platform
    ysctl cloudapp softwarelist --platform=ubuntu
`)

func NewSoftwareListOptions(ioStreams clientcmd.IOStreams) *SoftwareListOptions {
	return &SoftwareListOptions{
		IOStreams:  ioStreams,
		PageOffset: 0,
		PageSize:   100,
		Platform:   "WINDOWS",
	}
}

func NewSoftwareList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSoftwareListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "softwarelist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List softwares",
		TraverseChildren:      true,
		Long:                  "List softwares",
		Example:               softwareListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().IntVar(&o.PageOffset, "pageOffset", o.PageOffset, "Page offset")
	cmd.Flags().IntVar(&o.PageSize, "pageSize", o.PageSize, "Page size")
	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Software name")
	cmd.Flags().StringVar(&o.Platform, "platform", o.Platform, "Software platform")
	cmd.Flags().StringVar(&o.Zone, "zone", o.Zone, "Software zone")

	return cmd
}

func (o *SoftwareListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *SoftwareListOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Software.User.List(
		o.hc.API.CloudApp.Software.User.List.PageOffset(o.PageOffset),
		o.hc.API.CloudApp.Software.User.List.PageSize(o.PageSize),
		o.hc.API.CloudApp.Software.User.List.Name(o.Name),
		o.hc.API.CloudApp.Software.User.List.Platform(o.Platform),
		o.hc.API.CloudApp.Software.User.List.Zone(o.Zone),
	)

	if err != nil {
		return err
	}

	for _, software := range r.Data.Software {
		bf := bytes.NewBuffer([]byte{})
		jsonEncode := json.NewEncoder(bf)
		jsonEncode.SetEscapeHTML(false)
		if err := jsonEncode.Encode(software); err != nil {
			return err
		}

		fmt.Fprintf(o.Out, "%s\n", bf.String())
	}
	return nil
}
