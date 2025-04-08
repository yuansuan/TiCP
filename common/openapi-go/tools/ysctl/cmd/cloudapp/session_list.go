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

type SessionListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	SessionIds, Status, Zone string
	PageOffset, PageSize     int
}

var sessionListExample = templates.Examples(`
    # List all sessions
    ysctl cloudapp sessionlist

    # List all sessions with page offset and page size
    ysctl cloudapp sessionlist --pageOffset=0 --pageSize=100

    # List all sessions with session ids
    ysctl cloudapp sessionlist --sessionIDs=4WoY5JA8mvE

    # List all sessions with status
    ysctl cloudapp sessionlist --status=running
`)

func NewSessionListOptions(ioStreams clientcmd.IOStreams) *SessionListOptions {
	return &SessionListOptions{
		IOStreams:  ioStreams,
		PageOffset: 0,
		PageSize:   100,
	}
}

func NewSessionList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewSessionListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sessionlist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List sessions",
		TraverseChildren:      true,
		Long:                  "List sessions",
		Example:               sessionListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.SessionIds, "sessionIDs", o.SessionIds, "sessionIDs")
	cmd.Flags().StringVar(&o.Status, "status", o.Status, "status")
	cmd.Flags().StringVar(&o.Zone, "zone", o.Zone, "zone")
	cmd.Flags().IntVar(&o.PageOffset, "pageOffset", o.PageOffset, "pageOffset")
	cmd.Flags().IntVar(&o.PageSize, "pageSize", o.PageSize, "pageSize")

	return cmd
}

func (o *SessionListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *SessionListOptions) Run(args []string) error {
	r, err := o.hc.API.CloudApp.Session.User.List(
		o.hc.API.CloudApp.Session.User.List.PageOffset(o.PageOffset),
		o.hc.API.CloudApp.Session.User.List.PageSize(o.PageSize),
		o.hc.API.CloudApp.Session.User.List.SessionIds(o.SessionIds),
		o.hc.API.CloudApp.Session.User.List.Status(o.Status),
		o.hc.API.CloudApp.Session.User.List.Zone(o.Zone),
	)

	if err != nil {
		return err
	}

	for _, item := range r.Data.Sessions {
		bf := bytes.NewBuffer([]byte{})
		jsonEncode := json.NewEncoder(bf)
		jsonEncode.SetEscapeHTML(false)
		if err := jsonEncode.Encode(item); err != nil {
			return err
		}

		fmt.Fprintf(o.Out, "%s\n", bf.String())
	}
	return nil
}
