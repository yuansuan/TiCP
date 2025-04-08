package account

import (
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

type UserAddOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	Phone                   string
	Name                    string
	CompanyName             string
	UserChannel             string
	Password                string
	UnifiedSocialCreditCode string
	Email                   string
}

var addExample = templates.Examples(`
    # add a new ys user
    ysctl account adduser --phone=1999999999 --name=张三 --company_name=张三的公司 --user_channel=channel1 --unified_code=93234234324XCVQEFFE 
`)

func NewUserAddOptions(ioStreams clientcmd.IOStreams) *UserAddOptions {
	return &UserAddOptions{
		IOStreams: ioStreams,
	}
}

func NewUserAdd(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewUserAddOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "useradd",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "add user",
		TraverseChildren:      true,
		Long:                  "add user",
		Example:               addExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Phone, "phone", "", "phone num")
	cmd.Flags().StringVar(&o.Name, "name", "", "user name")
	cmd.Flags().StringVar(&o.CompanyName, "company_name", "", "company name")
	cmd.Flags().StringVar(&o.UserChannel, "user_channel", "", "user channel")
	cmd.Flags().StringVar(&o.Password, "password", "",
		"password, if empty random password will be generated")
	cmd.Flags().StringVar(&o.UnifiedSocialCreditCode, "unified_code", "", "企业信用代码")
	cmd.Flags().StringVar(&o.Email, "email", "", "email")
	return cmd
}

func (o *UserAddOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, err = f.RESTClient()
	if err != nil {
		return err
	}
	return nil
}
