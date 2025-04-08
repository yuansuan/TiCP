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

type SessionRebootOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionRebootExample = templates.Examples(`
    # Reboot session
    ysctl cloudapp sessionReboot --sessionID=4WoY5JA8mvE
`)

func NewSessionRebootOptions(ioStreams clientcmd.IOStreams) *SessionRebootOptions {
	return &SessionRebootOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionReboot(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionRebootOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionReboot",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Reboot session",
		TraverseChildren:      true,
		Long:                  "Reboot session",
		Example:               sessionRebootExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionRebootOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionRebootOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.Reboot(
		o.hc.API.CloudApp.Session.User.Reboot.Id(o.SessionID),
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
