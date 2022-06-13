package auth

import (
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func displayLogout(r *config.Red_t, cmd *cobra.Command) {
	r.SetApiKey("")
	r.SetServer("")
	if err := r.Save(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("you have successfuly logged out")
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
