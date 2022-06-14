package auth

import (
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

const (
	FLAG_SERVER   = "server"
	FLAG_USERNAME = "username"
	FLAG_APIKEY   = "apikey"
)

func NewCmdAuth(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth to Redmine",
		Long:  "Authenticate to Redmine server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println("Could not get usage info, abort..")
			}
		},
	}

	cmd.AddCommand(cmdAuthLogin(r))
	cmd.AddCommand(cmdAuthLogout(r))

	cmd.PersistentFlags().String(FLAG_SERVER, "", "URL to redmine server")
	cmd.PersistentFlags().String(FLAG_USERNAME, "", "Username to redmine")
	cmd.PersistentFlags().String(FLAG_APIKEY, "", "Use ApiKey instead of username and password")

	return cmd
}
