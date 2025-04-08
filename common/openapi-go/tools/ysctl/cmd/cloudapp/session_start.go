package cloudapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cloudappstart "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/start"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type SessionStartOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	HardwareID, SoftwareID string
	MountPaths             string
	PeriodType             string
	PeriodNum              int
	ChargeType             string
	PayByAccessKeyID       string
	PayByAccessSecret      string
}

var sessionStartExample = templates.Examples(`
    # Start session
    ysctl cloudapp sessionstart --hardwareID=4WoY5JA8mvE --softwareID=4WoY5JA8mvE --mountPaths=/data=/data,/bin=/bin --periodType=hour --periodNum=1 --chargeType=PrePaid
`)

func NewSessionStartOptions(ioStreams clientcmd.IOStreams) *SessionStartOptions {
	return &SessionStartOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionStart(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionStartOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionstart",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Start session",
		TraverseChildren:      true,
		Long:                  "Start session",
		Example:               sessionStartExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.HardwareID, "hardwareID", "", "Hardware ID")
	cmd.Flags().StringVar(&o.SoftwareID, "softwareID", "", "Software ID")
	cmd.Flags().StringVar(&o.MountPaths, "mountPaths", "", "Mount paths")
	cmd.Flags().StringVar(&o.PeriodType, "periodType", "", "Period type")
	cmd.Flags().IntVar(&o.PeriodNum, "periodNum", 0, "Period num")
	cmd.Flags().StringVar(&o.ChargeType, "chargeType", "PostPaid", "Charge type: PrePaid or PostPaid")

	return cmd
}

func (o *SessionStartOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *SessionStartOptions) Run(args []string) error {
	mountPaths := make(map[string]string)

	if o.MountPaths != "" {
		s := strings.Split(o.MountPaths, ",")
		for _, v := range s {
			m := strings.Split(v, "=")
			mountPaths[m[0]] = m[1]
		}
	}
	t := v20230530.ChargeType(o.ChargeType)

	chargeParams := v20230530.ChargeParams{
		ChargeType: &t,
		PeriodType: &o.PeriodType,
		PeriodNum:  &o.PeriodNum,
	}

	opts := make([]cloudappstart.Option, 0)
	opts = append(opts, o.hc.API.CloudApp.Session.User.Start.HardwareId(o.HardwareID),
		o.hc.API.CloudApp.Session.User.Start.SoftwareId(o.SoftwareID),
		o.hc.API.CloudApp.Session.User.Start.ChargeParams(chargeParams),
	)
	if mountPaths != nil {
		opts = append(opts, o.hc.API.CloudApp.Session.User.Start.MountPaths(mountPaths))
	}
	if o.PayByAccessKeyID != "" && o.PayByAccessSecret != "" {
		opts = append(opts, o.hc.API.CloudApp.Session.User.Start.PayByParams(o.PayByAccessKeyID, o.PayByAccessSecret))
	}
	r, err := o.hc.API.CloudApp.Session.User.Start(opts...)

	if err != nil {
		return err
	}

	bf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(r.Data); err != nil {
		return err
	}
	fmt.Println(bf.String())
	return nil
}
