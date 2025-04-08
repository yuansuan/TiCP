package account

import (
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type ExchangeOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	Phone    string
	Email    string
	YsId     string
	Password string
}

var exchangeExample = templates.Examples(`
    # add a new ys user
    ysctl account exchange_ak --phone=1999999999 --email=abc@yuansuan.cn --ys_id=XDFAF344X --password=password123
`)

func NewExchangeOptions(ioStreams clientcmd.IOStreams) *ExchangeOptions {
	return &ExchangeOptions{
		IOStreams: ioStreams,
	}
}

func NewExchange(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewExchangeOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "exchange_ak",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "exchange access key",
		TraverseChildren:      true,
		Long:                  "exchange access key",
		Example:               addExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Phone, "phone", "", "phone num")
	cmd.Flags().StringVar(&o.Email, "email", "", "email")
	cmd.Flags().StringVar(&o.YsId, "ys_id", "", "ys id")
	cmd.Flags().StringVar(&o.Password, "password", "", "password")
	return cmd
}

func (o *ExchangeOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}
