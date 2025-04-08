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

type HardwarGetOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	HardwareID string
}

var hardwareGetExample = templates.Examples(`
    # Get hardware
    ysctl cloudapp hardwareget --hardwareID=4WoY5JA8mvE
`)

func NewHardwareGetOptions(ioStreams clientcmd.IOStreams) *HardwarGetOptions {
	return &HardwarGetOptions{
		IOStreams: ioStreams,
	}
}

func NewHardwareGet(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewHardwareGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "hardwareget",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Get hardware",
		TraverseChildren:      true,
		Long:                  "Get hardware",
		Example:               hardwareGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.HardwareID, "hardwareID", o.HardwareID, "Hardware ID")

	return cmd
}

func (o *HardwarGetOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *HardwarGetOptions) Run(args []string) error {
	hardware, err := o.hc.API.CloudApp.Hardware.User.Get(
		o.hc.API.CloudApp.Hardware.User.Get.Id(o.HardwareID),
	)
	if err != nil {
		return err
	}
	bf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(hardware); err != nil {
		return err
	}
	fmt.Println(bf.String())
	return nil
}
