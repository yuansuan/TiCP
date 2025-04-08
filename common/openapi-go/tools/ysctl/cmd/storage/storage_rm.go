package storage

import (
	"fmt"

	openys "github.com/yuansuan/ticp/common/openapi-go"
	rm "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/rm"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"

	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"

	"github.com/spf13/cobra"
)

type StorageRmOptions struct {
	Path           string
	IgnoreNotExist bool
	hc             *openys.Client
	clientcmd.IOStreams
}

var rmExample = templates.Examples(`
	# Remove a file
	ysctl storage rm --path=/4TiSsZonTa3/file.txt

	# Remove a file ignoring non-existence
	ysctl storage rm --path=/4TiSsZonTa3/file.txt --ignore-not-exist
`)

func NewStorageRmOptions(ioStreams clientcmd.IOStreams) *StorageRmOptions {
	return &StorageRmOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageRm(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageRmOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "rm",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Remove files or directories",
		TraverseChildren:      true,
		Long:                  "Remove files or directories",
		Example:               rmExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.Path, "path", "", "Path to the file or directory to remove")
	cmd.Flags().BoolVar(&o.IgnoreNotExist, "ignore-not-exist", false, "Ignore if the file or directory does not exist")
	cmd.MarkFlagRequired("path")
	return cmd
}

func (o *StorageRmOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, _, err = f.StorageClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *StorageRmOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "Removing: %s\n", o.Path)
	options := []rm.Option{
		o.hc.Storage.Rm.Path(o.Path),
		o.hc.Storage.Rm.IgnoreNotExist(o.IgnoreNotExist),
	}

	_, err := o.hc.API.Storage.Rm(options...)
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Successfully removed: %s\n", o.Path)
	return nil
}
