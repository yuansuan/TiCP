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

type SessionPowerOnOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionPowerOnExample = templates.Examples(`
    # PowerOn session
    ysctl cloudapp sessionPowerOn --sessionID=4WoY5JA8mvE
`)

func NewSessionPowerOnOptions(ioStreams clientcmd.IOStreams) *SessionPowerOnOptions {
	return &SessionPowerOnOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionPowerOn(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionPowerOnOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionPowerOn",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "PowerOn session",
		TraverseChildren:      true,
		Long:                  "PowerOn session",
		Example:               sessionPowerOnExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionPowerOnOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionPowerOnOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.PowerOn(
		o.hc.API.CloudApp.Session.User.PowerOn.Id(o.SessionID),
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
