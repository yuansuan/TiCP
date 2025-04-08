package singularity

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
)

type PushCommand *cobra.Command

func NewPushCommand(cli registry.Client) PushCommand {
	cmd := &cobra.Command{
		Use:   "push [flags] <image-locator>",
		Short: "push locally image or file to remote",
		Args:  cobra.ExactArgs(1),
	}

	var options struct {
		Force     bool
		ImageFile string
		Alias     string
	}

	cmd.Flags().BoolVar(&options.Force, "force", false, "force to overwrite remotely image")
	cmd.Flags().StringVarP(&options.ImageFile, "file", "f", "", "the image to upload")
	cmd.Flags().StringVarP(&options.ImageFile, "alias", "l", "", "upload locatable image")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(options.ImageFile) != 0 && len(options.Alias) != 0 {
			return errors.New("only one file or locator can be set")
		}

		locator, err := image.FromString(args[0])
		if err != nil {
			return err
		}

		var opts []registry.PushOption
		opts = append(opts, registry.WithPushForce(options.Force))
		opts = append(opts, registry.WithPushLogger(LogToConsole))
		opts = append(opts, registry.WithPushContext(cmd.Context()))
		opts = append(opts, registry.WithPushProgressBar(GetProgressBarOpts(locator.ShortString())))
		if len(options.ImageFile) != 0 {
			opts = append(opts, registry.WithPushLocalFile(options.ImageFile))
		} else if len(options.Alias) != 0 {
			l, err := image.FromString(options.Alias)
			if err != nil {
				return errors.Wrap(err, "alias")
			}

			opts = append(opts, registry.WithPushLocalAlias(l))
		}

		remotely, err := cli.Push(locator, opts...)
		if err != nil {
			return err
		}

		Clog("Pushed", remotely.Locator().String(), nil)
		return nil
	}

	return cmd
}
