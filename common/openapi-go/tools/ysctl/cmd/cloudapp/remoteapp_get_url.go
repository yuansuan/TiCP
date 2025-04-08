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

type RemoteAppOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionId, RemoteAppName string
}

var remoteAppExample = templates.Examples(`
    # Get remote app url
    ysctl cloudapp remoteappgeturl --sessionId=4WoY5JA8mvE --remoteAppName=ys-1c1g1g1s
`)

func NewRemoteAppOptions(ioStreams clientcmd.IOStreams) *RemoteAppOptions {
	return &RemoteAppOptions{
		IOStreams: ioStreams,
	}
}

func NewRemoteApp(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewRemoteAppOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "remoteappgeturl",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Get remote app url",
		TraverseChildren:      true,
		Long:                  "Get remote app url",
		Example:               remoteAppExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.SessionId, "sessionId", o.SessionId, "Session ID")
	cmd.Flags().StringVar(&o.RemoteAppName, "remoteAppName", o.RemoteAppName, "Remote app name")

	return cmd
}

func (o *RemoteAppOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *RemoteAppOptions) Run(args []string) error {
	r, err := o.hc.CloudApp.RemoteApp.User.Get(
		o.hc.CloudApp.RemoteApp.User.Get.RemoteAppName(o.RemoteAppName),
		o.hc.CloudApp.RemoteApp.User.Get.SessionId(o.SessionId),
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
