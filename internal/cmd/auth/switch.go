package auth

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func cmdAuthSwitch(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch",
		Short: "switch current authenticated Redmine server",
		Long:  "Switch the current authenticated Redmine server by name",
		Run: func(cmd *cobra.Command, args []string) {
			servers := r.GetServers()
			if servers == nil {
				print.Info("No server is currently authenticated.")
				return
			}

			var names []string
			for _, server := range servers {
				names = append(names, server.Name)
			}

			name, serverID := r.Term.ChooseString("Select server to switch to", names)

			r.SetDefaultServerById(serverID)

			if err := r.Save(); err != nil {
				print.Error("Could not save configuration: %s", err.Error())
				return
			}
			print.OK("Switched to server '%s' successfully.", name)
		},
	}
	return cmd
}
