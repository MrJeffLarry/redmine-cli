package auth

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func displayLogout(r *config.Red_t, cmd *cobra.Command) {
	r.ClearAll()
	if err := r.Save(); err != nil {
		print.Error(err.Error())
		return
	}
	print.OK("you have successfully logged out")
}

func cmdAuthLogout(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "logout from Redmine",
		Long:  "Rest and logout from Redmine server",
		Run: func(cmd *cobra.Command, args []string) {
			displayLogout(r, cmd)
		},
	}
	return cmd
}
