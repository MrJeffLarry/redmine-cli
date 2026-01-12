package auth

import (
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func displayList(r *config.Red_t, cmd *cobra.Command) {
	servers := r.GetServers()
	if servers == nil {
		print.Info("No server is currently authenticated.")
		return
	}

	head := []string{"ID", "Active", "Name", "Url"}

	l := print.NewList(head...)

	for serverId, server := range servers {
		id := print.Column{}
		def := print.Column{}
		name := print.Column{}
		url := print.Column{}

		id.Content = strconv.FormatInt(int64(serverId), 10)
		id.FgColor = print.ID

		if serverId == r.Config.DefaultServer {
			def.Content = "*"
			def.FgColor = print.Green
		} else {
			def.Content = ""
		}

		name.Content = server.Name
		url.Content = server.Server

		l.AddRow(serverId, -1, id, def, name, url)
	}
	l.SetTotal(len(servers))

	editor.StartPage(r.Config.Pager, l.Render())
}

func cmdAuthList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "display list of authenticated Redmine servers",
		Long:  "Display list of authenticated Redmine servers",
		Run: func(cmd *cobra.Command, args []string) {
			displayList(r, cmd)
		},
	}
	return cmd
}
