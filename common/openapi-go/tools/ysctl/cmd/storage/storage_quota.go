package storage

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type StorageQuotaOptions struct {
	hc     *openys.Client
	config *clientcmd.Config
	clientcmd.IOStreams
}

var quotaExample = templates.Examples(`

	# get storage usage and quota limit of current user,unit is GB
	ysctl storage quota
`)

func NewStorageQuotaOptions(ioStreams clientcmd.IOStreams) *StorageQuotaOptions {
	return &StorageQuotaOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageQuota(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageQuotaOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "quota",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "get storage quota info",
		TraverseChildren:      true,
		Long:                  "get storage quota info",
		Example:               quotaExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	return cmd
}

func (o *StorageQuotaOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, o.config, err = f.StorageClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *StorageQuotaOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "current user ak: %s,storage endpoint: %s\n", o.config.AccessKeyID, o.config.StorageEndpoint)
	r, err := o.hc.API.StorageQuota.GetQuotaAPI()
	if err != nil {
		return err
	}

	res, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	fmt.Fprintf(o.Out, "%s\n", res)

	return nil
}
