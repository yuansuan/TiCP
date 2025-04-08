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

type SessionPowerOffOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionPowerOffExample = templates.Examples(`
    # PowerOff session
    ysctl cloudapp sessionPowerOff --sessionID=4WoY5JA8mvE
`)

func NewSessionPowerOffOptions(ioStreams clientcmd.IOStreams) *SessionPowerOffOptions {
	return &SessionPowerOffOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionPowerOff(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionPowerOffOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionPowerOff",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "PowerOff session",
		TraverseChildren:      true,
		Long:                  "PowerOff session",
		Example:               sessionPowerOffExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionPowerOffOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionPowerOffOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.PowerOff(
		o.hc.API.CloudApp.Session.User.PowerOff.Id(o.SessionID),
	)
	if err != nil {
		return err
	}

	bf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(r); err != nil {
		return err
	}
	fmt.Println(bf.String())
	return nil
}
