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

type SessionGetOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionID string
}

var sessionGetExample = templates.Examples(`
    # Get session
    ysctl cloudapp sessionGet --sessionID=4WoY5JA8mvE
`)

func NewSessionGetOptions(ioStreams clientcmd.IOStreams) *SessionGetOptions {
	return &SessionGetOptions{
		IOStreams: ioStreams,
	}
}

func NewSessionGet(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionGet",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Get session",
		TraverseChildren:      true,
		Long:                  "Get session",
		Example:               sessionGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionID, "sessionID", o.SessionID, "Session ID")

	return cmd
}

func (o *SessionGetOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *SessionGetOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.Get(
		o.hc.API.CloudApp.Session.User.Get.Id(o.SessionID),
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
