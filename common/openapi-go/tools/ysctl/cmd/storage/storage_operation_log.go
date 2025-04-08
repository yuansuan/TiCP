package storage

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
	"time"
)

type StorageOperationLogOptions struct {
	hc            *openys.Client
	PageOffset    int64
	PageSize      int64
	FileName      string
	FileType      string
	OperationType string
	BeginTime     string
	EndTime       string
	config        *clientcmd.Config
	clientcmd.IOStreams
}

var operationLogExample = templates.Examples(`

	# Display all file upload type operation log with offset 0 and limit 1000 where file name contains "sim" and operation time between 2023-11-13 13:30:08 and 2023-11-13 13:32:30
	ysctl storage operation-log -o 0 -s 1000 -n sim -F FILE -O UPLOAD -b "2023-11-13 13:30:08" -e "2023-11-13 13:32:30"
`)

func NewStorageOperationLogOptions(ioStreams clientcmd.IOStreams) *StorageOperationLogOptions {
	return &StorageOperationLogOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageOperationLog(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageOperationLogOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "operation-log",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "list storage operation log",
		TraverseChildren:      true,
		Long:                  "list storage operation log",
		Example:               operationLogExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().Int64VarP(&o.PageOffset, "pageOffset", "o", 0, "Page offset")
	cmd.Flags().Int64VarP(&o.PageSize, "pageSize", "s", 100, "Page size")
	cmd.Flags().StringVarP(&o.FileName, "fileName", "n", "", "File name")
	cmd.Flags().StringVarP(&o.FileType, "fileType", "F", "", "File type, options: FILE, FOLDER, BATCH")
	cmd.Flags().StringVarP(&o.OperationType, "operationType", "O", "", "Operation type, options: UPLOAD, DOWNLOAD, DELETE, MOVE, MKDIR, COPY, COPY_RANGE,COMPRESS, CREATE, LINK, READ_AT, WRITE_AT")
	cmd.Flags().StringVarP(&o.BeginTime, "beginTime", "b", "", "Begin time format: YYYY-MM-DD HH:mm:ss")
	cmd.Flags().StringVarP(&o.EndTime, "endTime", "e", "", "End time format: YYYY-MM-DD HH:mm:ss")

	return cmd
}

func (o *StorageOperationLogOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, o.config, err = f.StorageClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *StorageOperationLogOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "current user ak: %s,storage endpoint: %s\n", o.config.AccessKeyID, o.config.StorageEndpoint)
	var beginTimestamp int64
	var endTimestamp int64
	var err error
	if o.BeginTime != "" {
		beginTimestamp, err = stringToUnixTimestamp(o.BeginTime)
		if err != nil {
			fmt.Fprintf(o.Out, "BeginTimeError,format: YYYY-MM-DD HH:mm:ss, Error: %s\n", err.Error())
			return err
		}
	}

	if o.EndTime != "" {
		endTimestamp, err = stringToUnixTimestamp(o.EndTime)
		if err != nil {
			fmt.Fprintf(o.Out, "EndTimeError,format: YYYY-MM-DD HH:mm:ss, Error: %s\n", err.Error())
			return err
		}
	}
	r, err := o.hc.API.StorageOperationLog.ListOperationLogAPI(
		o.hc.API.StorageOperationLog.ListOperationLogAPI.PageOffset(o.PageOffset),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.PageSize(o.PageSize),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.FileName(o.FileName),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.FileTypes(o.FileType),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.OperationTypes(o.OperationType),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.BeginTime(beginTimestamp),
		o.hc.API.StorageOperationLog.ListOperationLogAPI.EndTime(endTimestamp),
	)
	if err != nil {
		return err
	}

	data := make([][]string, 0, len(r.Data.OperationLog))
	table := tablewriter.NewWriter(o.Out)

	for _, f := range r.Data.OperationLog {
		srcPath := f.SrcPath
		if srcPath == "" {
			srcPath = "  "
		}

		destPath := f.DestPath
		if destPath == "" {
			destPath = "  "
		}

		data = append(data, []string{
			f.FileName,
			f.SrcPath,
			f.DestPath,
			f.FileType,
			f.OperationType,
			f.Size,
			formatTime(f.CreateTime),
		})
	}

	table = setOperationLogHeader(table)
	table = operationLogTableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func operationLogTableWriterDefaultConfig(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	return table
}

func setOperationLogHeader(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{"FileName", "SrcPath", "DestPath", "FileType", "OperationType", "Size", "CreateTime"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
	)

	return table
}

func stringToUnixTimestamp(str string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
