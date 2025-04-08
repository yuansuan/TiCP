package singularity

import (
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
)

type PullCommand *cobra.Command

// NewPullCommand 镜像拉取命令
func NewPullCommand(cli registry.Client) PullCommand {
	cmd := &cobra.Command{
		Use:          "pull [flags] images...",
		Short:        "pull image from remote",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: false,
	}

	var options struct {
		Force    bool
		Blocking bool
	}

	cmd.Flags().BoolVar(&options.Force, "force", false, "force to overwrite locally image")
	cmd.Flags().BoolVarP(&options.Blocking, "blocking", "b", false, "blocking until pulled")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			locator, err := image.FromString(arg)
			if err != nil {
				return err
			}

			var opts []registry.PullOption
			opts = append(opts, registry.WithPullForce(options.Force))
			opts = append(opts, registry.WithPullContext(cmd.Context()))
			opts = append(opts, registry.WithPullLogger(LogToConsole))
			opts = append(opts, registry.WithPullProgressBar(GetProgressBarOpts(locator.ShortString())))
			if options.Blocking {
				opts = append(opts, registry.WithPullBlocking())
			}

			locally, err := cli.Pull(cmd.Context(), locator, opts...)
			if err != nil {
				return err
			}

			Clog("Pulled", locally.Locator().String(), nil)
		}
		return nil
	}

	return cmd
}
