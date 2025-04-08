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

type SessionCloseOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionCloseExample = templates.Examples(`
    # Close session
    ysctl cloudapp sessionClose --sessionID=4WoY5JA8mvE
`)

func NewSessionCloseOptions(ioStreams clientcmd.IOStreams) *SessionCloseOptions {
	return &SessionCloseOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionClose(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionCloseOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionClose",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Close session",
		TraverseChildren:      true,
		Long:                  "Close session",
		Example:               sessionCloseExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionCloseOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionCloseOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.Close(
		o.hc.API.CloudApp.Session.User.Close.Id(o.SessionID),
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
