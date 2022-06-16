package user

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdUser(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get users info",
		Long:  "Gets information about users account",
		Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdUserMe(r))

	return cmd
}
