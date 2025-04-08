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

type SessionReadyOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionReadyExample = templates.Examples(`
    # Ready session
    ysctl cloudapp sessionready --sessionID=4WoY5JA8mvE
`)

func NewSessionReadyOptions(ioStreams clientcmd.IOStreams) *SessionReadyOptions {
	return &SessionReadyOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionReady(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionReadyOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionready",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Ready session",
		TraverseChildren:      true,
		Long:                  "Ready session",
		Example:               sessionReadyExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionReadyOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionReadyOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.Ready(
		o.hc.API.CloudApp.Session.User.Ready.Id(o.SessionID),
	)
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
