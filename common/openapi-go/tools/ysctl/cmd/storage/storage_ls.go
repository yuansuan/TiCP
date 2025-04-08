package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type StorageListOptions struct {
	PageOffset   int64
	PageSize     int64
	Path         string
	FilterRegexp string
	hc           *openys.Client
	clientcmd.IOStreams
}

var listExample = templates.Examples(`
	# Display some files
	ysctl storage ls --path=/4TiSsZonTa3 --filterRegexp=^/test/.*$  --pageOffset=0 --pageSize=10
`)

func NewStorageListOptions(ioStreams clientcmd.IOStreams) *StorageListOptions {
	return &StorageListOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageList(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "ls",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "List files",
		TraverseChildren:      true,
		Long:                  "List files",
		Example:               listExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().Int64Var(&o.PageOffset, "pageOffset", 0, "Page offset")
	cmd.Flags().Int64Var(&o.PageSize, "pageSize", 100, "Page size")
	cmd.Flags().StringVar(&o.Path, "path", "", "Path")
	cmd.Flags().StringVar(&o.FilterRegexp, "filterRegexp", "", "Filter regexp")
	return cmd
}

func (o *StorageListOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, _, err = f.StorageClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *StorageListOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "Current Config: PageOffset: %d, PageSize: %d, Path: %s, FilterRegexp: %s\n", o.PageOffset, o.PageSize, o.Path, o.FilterRegexp)
	r, err := o.hc.API.Storage.LsWithPage(
		o.hc.Storage.LsWithPage.Path(o.Path),
		o.hc.Storage.LsWithPage.PageOffset(o.PageOffset),
		o.hc.Storage.LsWithPage.PageSize(o.PageSize),
		o.hc.Storage.LsWithPage.FilterRegexp(o.FilterRegexp),
	)
	if err != nil {
		return err
	}

	data := make([][]string, 0, len(r.Data.Files))
	table := tablewriter.NewWriter(o.Out)

	for _, f := range r.Data.Files {
		data = append(data, []string{
			f.Name,
			strconv.FormatInt(f.Size, 10),
			strconv.FormatUint(uint64(f.Mode), 8),
			formatTime(f.ModTime),
			strconv.FormatBool(f.IsDir),
		})
	}

	table = setHeader(table)
	table = TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
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
	table.SetHeader([]string{"Name", "Size", "Mode", "ModTime", "IsDir"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgGreenColor})

	return table
}
