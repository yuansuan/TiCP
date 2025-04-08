package cloudapp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/api/list"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type HardwareListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	PageOffset, PageSize int
	Name, Zone           string
	Cpu, Mem, Gpu        int
}

var hardwareListExample = templates.Examples(`
    # List all hardwares
    ysctl cloudapp hardwarelist

    # List all hardwares with page offset and page size
    ysctl cloudapp hardwarelist --pageOffset=0 --pageSize=100

    # List all hardwares with name
    ysctl cloudapp hardwarelist --name=ys-1c1g1g1s

    # List all hardwares with cpu
    ysctl cloudapp hardwarelist --cpu=1
`)

func NewHardwareListOptions(ioStreams clientcmd.IOStreams) *HardwareListOptions {
	return &HardwareListOptions{
		IOStreams:  ioStreams,
		PageOffset: 0,
		PageSize:   100,
	}
}

func NewHardwareList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewHardwareListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "hardwarelist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List hardwares",
		TraverseChildren:      true,
		Long:                  "List hardwares",
		Example:               hardwareListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().IntVar(&o.PageOffset, "pageOffset", o.PageOffset, "Page offset")
	cmd.Flags().IntVar(&o.PageSize, "pageSize", o.PageSize, "Page size")
	cmd.Flags().IntVar(&o.Cpu, "cpu", o.Cpu, "Cpu")
	cmd.Flags().IntVar(&o.Mem, "mem", o.Mem, "Mem")
	cmd.Flags().IntVar(&o.Gpu, "gpu", o.Gpu, "Gpu")
	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Name")
	cmd.Flags().StringVar(&o.Zone, "zone", o.Zone, "Zone")

	return cmd
}

func (o *HardwareListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *HardwareListOptions) Run(args []string) error {
	opts := o.ensureHardwareListOpts()
	opts = append(opts,
		o.hc.API.CloudApp.Hardware.User.List.PageOffset(o.PageOffset),
		o.hc.API.CloudApp.Hardware.User.List.PageSize(o.PageSize))

	r, err := o.hc.API.CloudApp.Hardware.User.List(opts...)
	if err != nil {
		return err
	}

	for _, hardware := range r.Data.Hardware {
		bf := bytes.NewBuffer([]byte{})
		jsonEncode := json.NewEncoder(bf)
		jsonEncode.SetEscapeHTML(false)
		if err := jsonEncode.Encode(hardware); err != nil {
			return err
		}

		fmt.Fprintf(o.Out, "%s\n", bf.String())
	}
	return nil
}

func (o *HardwareListOptions) ensureHardwareListOpts() []list.Option {
	opts := make([]list.Option, 0)
	if o.Name != "" {
		opts = append(opts, o.hc.API.CloudApp.Hardware.User.List.Name(o.Name))
	}
	if o.Cpu > 0 {
		opts = append(opts, o.hc.API.CloudApp.Hardware.User.List.Cpu(o.Cpu))
	}
	if o.Mem > 0 {
		opts = append(opts, o.hc.API.CloudApp.Hardware.User.List.Mem(o.Mem))
	}
	if o.Gpu >= 0 {
		opts = append(opts, o.hc.API.CloudApp.Hardware.User.List.Gpu(o.Gpu))
	}
	if o.Zone != "" {
		opts = append(opts, o.hc.API.CloudApp.Hardware.User.List.Zone(o.Zone))
	}

	return opts
}
