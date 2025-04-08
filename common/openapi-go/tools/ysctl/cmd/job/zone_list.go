package job

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type ZoneListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
}

var zoneListExample = templates.Examples(`
    # List all zones
    ysctl job zonelist
`)

func NewZoneListOptions(ioStreams clientcmd.IOStreams) *ZoneListOptions {
	return &ZoneListOptions{
		IOStreams: ioStreams,
	}
}

func NewZoneList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewZoneListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "zonelist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List zones",
		TraverseChildren:      true,
		Long:                  "List zones",
		Example:               zoneListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	return cmd
}

func (o *ZoneListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *ZoneListOptions) Run(args []string) error {
	r, err := o.hc.API.Job.ZoneList()
	if err != nil {
		return err
	}

	data := make([][]string, len(r.Data.Zones))
	table := tablewriter.NewWriter(o.Out)

	for zone := range r.Data.Zones {
		data = append(data, []string{
			zone,
			r.Data.Zones[zone].StorageEndpoint,
		})
	}
	table = setZoneHeader(table)
	table = TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()
	return nil
}

func setZoneHeader(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{
		"Zone",
		"StorageEndpoint",
	})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor})
	return table
}
