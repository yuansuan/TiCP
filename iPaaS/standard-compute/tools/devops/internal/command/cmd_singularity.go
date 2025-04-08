package command

import (
	"github.com/spf13/cobra"

	sg "github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops/internal/command/singularity"
)

type SingularityCommand *cobra.Command

func NewSingularityCommand(pull sg.PullCommand, push sg.PushCommand, search sg.SearchCommand) SingularityCommand {
	cmd := &cobra.Command{
		Use:     "singularity",
		Short:   "Managing local and remote images of singularity (alias: sg)",
		Aliases: []string{"sg"},
	}

	cmd.AddCommand(pull)
	cmd.AddCommand(push)
	cmd.AddCommand(search)

	return cmd
}
