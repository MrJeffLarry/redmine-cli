package users

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdUsers(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Get users info",
		Long:  "Gets information about users account",
		Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdUsersMe(r))

	return cmd
}
