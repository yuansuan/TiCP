package singularity

import (
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"
)

type SearchCommand *cobra.Command

func NewSearchCommand(cli registry.Client) SearchCommand {
	cmd := &cobra.Command{
		Use:   "search [flags] search-pattern",
		Short: "search the remote and local registry for images",
		Args:  cobra.MaximumNArgs(1),
	}

	var options struct {
		Defaults bool
		Locally  bool
	}

	cmd.Flags().BoolVarP(&options.Defaults, "defaults", "d", false, "search defaults image")
	cmd.Flags().BoolVarP(&options.Locally, "locally", "l", false, "search image from locally")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var pattern string
		if len(args) == 1 {
			pattern = args[0]
		}

		var opts []registry.SearchOption
		opts = append(opts, registry.WithSearchContext(cmd.Context()))
		opts = append(opts, registry.WithSearchDefaults(options.Defaults))
		opts = append(opts, registry.WithSearchLocally(options.Locally))

		locators, err := cli.Search(pattern, opts...)
		if err != nil {
			return err
		}

		PrintDefaultedLocators(locators)
		return nil
	}

	return cmd
}
