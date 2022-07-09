package auth

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func NewCmdAuth(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth to Redmine",
		Long:  "Authenticate to Redmine server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				print.Error("Could not get usage info, abort..")
			}
		},
	}

	cmd.AddCommand(cmdAuthLogin(r))
	cmd.AddCommand(cmdAuthLogout(r))

	return cmd
}
