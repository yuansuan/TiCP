package job

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type AppListOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
}

var listExample = templates.Examples(`
    # List all applications
    ysctl job applist
`)

func NewAppListOptions(ioStreams clientcmd.IOStreams) *AppListOptions {
	return &AppListOptions{
		IOStreams: ioStreams,
	}
}

func NewAppList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewAppListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "applist",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List applications",
		TraverseChildren:      true,
		Long:                  "List applications",
		Example:               listExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	return cmd
}

func (o *AppListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *AppListOptions) Run(args []string) error {
	r, err := o.hc.API.Job.ListAPP()
	if err != nil {
		return err
	}

	data := make([][]string, 0, len(r.Data))
	table := tablewriter.NewWriter(o.Out)

	for _, v := range r.Data {
		data = append(data, []string{
			v.AppID,
			v.Name,
			v.Type,
			v.Version,
			v.Image,
			v.Description,
			v.BinPath,
		})
	}

	table = setHeader(table)
	table = TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func TableWriterDefaultConfig(table *tablewriter.Table) *tablewriter.Table {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with two space
	table.SetNoWhiteSpace(true)

	return table
}

func setHeader(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{"id", "name", "type", "version", "image", "description", "bin_path"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
	)

	return table
}
