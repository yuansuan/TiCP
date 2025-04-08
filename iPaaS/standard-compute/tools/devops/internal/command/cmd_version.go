package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops"
)

type VersionCommand *cobra.Command

func NewVersionCommand(cfg *config.Singularity) VersionCommand {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the client and registry version information currently",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Version: v%s(Rev. %s)\n", devops.Version, devops.Revision)
			fmt.Printf("BuildTime: %s\n", devops.BuildTime)
			fmt.Printf("Registry: %s@%s\n", cfg.Registry.Region, cfg.Registry.Bucket)
		},
	}
}
