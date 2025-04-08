package cmd

import (
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/account"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/cloudapp"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/job"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/signature"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/storage"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
)

func NewDefaultYSCtlCommand() *cobra.Command {
	return NewYSCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewDefaultLiteYSCtlCommand() *cobra.Command {
	return NewLiteYSCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewYSCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ysctl",
		Short: "ysctl controls the OpenAPI",
		Long: templates.LongDesc(`
        ysctl controls the OpenAPI, is the client side tool for YSCloud.

		Config file is named ysctl.yaml, which is located in the current directory by default.

		There is four elements in the config file:
		- access_key_id: The access key id of the user.
		- access_key_secret: The access key secret of the user.
		- endpoint: The endpoint of the OpenAPI.
		- storage_endpoint: The endpoint of the storage (if you want to use storage commands).
        
        For more information, please contact PID.`),

		Run: runHelp,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			switch cmdutil.Loglevel {
			case "debug":
			case "info":
			case "warn":
			case "error":
			default:
				cmdutil.Loglevel = "info"
			}
			log, _ := logging.NewLogger(logging.WithDefaultLogConfigOption(),
				logging.WithLogLevel(logging.LogLevel(cmdutil.Loglevel)))
			*logging.Default() = *log
		},
	}

	cmd.PersistentFlags().StringVar(&cmdutil.Endpoint, "endpoint", "", "ys openapi endpoint")
	cmd.PersistentFlags().StringVar(&cmdutil.AccessKeyID, "access_id", "", "access key id")
	cmd.PersistentFlags().StringVar(&cmdutil.AccessKeySecret, "access_secret", "", "access key secret")
	cmd.PersistentFlags().StringVar(&cmdutil.Proxy, "proxy", "", "proxy to used")
	cmd.PersistentFlags().StringVar(&cmdutil.StorageEndpoint, "storage_endpoint", "", "ys storage endpoint")
	cmd.PersistentFlags().StringVar(&cmdutil.Loglevel, "log_level", "info", "log level, default info")

	cobra.OnInitialize(func() {
		loadConfig("ysctl")
	})

	ioStreams := clientcmd.IOStreams{In: in, Out: out, ErrOut: err}

	f := cmdutil.NewApiClient()

	groups := templates.CommandGroups{
		{
			Message: "Job commands:",
			Commands: []*cobra.Command{
				job.NewCmdJob(f, ioStreams),
			},
		},
		{
			Message: "Storage commands:",
			Commands: []*cobra.Command{
				storage.NewCmdStorage(f, ioStreams),
			},
		},
		{
			Message: "CloudApp commands:",
			Commands: []*cobra.Command{
				cloudapp.NewCmdCloudApp(f, ioStreams),
			},
		},
		{
			Message: "Account commands:",
			Commands: []*cobra.Command{
				account.NewCmdAccount(f, ioStreams),
			},
		},
		{
			Message: "Signature commands:",
			Commands: []*cobra.Command{
				signature.NewGenSigner(f, ioStreams),
			},
		},
	}
	groups.Add(cmd)
	return cmd
}

func NewLiteYSCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ysctl",
		Short: "ysctl controls the OpenAPI, lite version.",
		Long: templates.LongDesc(`
        ysctl controls the OpenAPI, is the client side tool for YSCloud, lite version.

		Config file is named ysctl.yaml, which is located in the current directory by default.

		There is four elements in the config file:
		- access_key_id: The access key id of the user.
		- access_key_secret: The access key secret of the user.
		- endpoint: The endpoint of the OpenAPI.
		- storage_endpoint: The endpoint of the storage (if you want to use storage commands).
        
        For more information, please contact PID.`),

		Run: runHelp,
	}

	cmd.PersistentFlags().StringVar(&cmdutil.Endpoint, "endpoint", "", "ys openapi endpoint")
	cmd.PersistentFlags().StringVar(&cmdutil.AccessKeyID, "access_id", "", "access key id")
	cmd.PersistentFlags().StringVar(&cmdutil.AccessKeySecret, "access_secret", "", "access key secret")
	cmd.PersistentFlags().StringVar(&cmdutil.Proxy, "proxy", "", "proxy to used")
	cmd.PersistentFlags().StringVar(&cmdutil.StorageEndpoint, "storage_endpoint", "", "ys storage endpoint")

	cobra.OnInitialize(func() {
		loadConfig("ysctl")
	})

	ioStreams := clientcmd.IOStreams{In: in, Out: out, ErrOut: err}

	f := cmdutil.NewApiClient()

	groups := templates.CommandGroups{
		{
			Message: "Storage commands:",
			Commands: []*cobra.Command{
				storage.NewLiteCmdStorage(f, ioStreams),
			},
		},
	}
	groups.Add(cmd)
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func loadConfig(defaultName string) {
	viper.AddConfigPath(".")
	viper.SetConfigName(defaultName)

	viper.SetConfigType("yaml")

	if _, err := os.Stat(fmt.Sprintf("%s.yaml", defaultName)); os.IsNotExist(err) {
		return
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
