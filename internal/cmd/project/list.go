package project

import (
	"encoding/json"
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/spf13/cobra"
)

const (
	FLAG_QUERY       = "query"
	FLAG_QUERY_SHORT = "q"
)

func displayListGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int
	head := []string{"ID", "NAME"}

	projects := projects{}

	path += util.ParseFlags(cmd, 0, []string{"id", "name"})

	if query, _ := cmd.Flags().GetString(FLAG_QUERY); query != "" {
		path += "name=~" + query + "&"
	}

	print.Debug(r, path)

	if body, status, err = api.ClientGET(r, path); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.Debug(r, "%d %s", string(body))

	if err := json.Unmarshal(body, &projects); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	l := print.NewList(head...)

	for _, project := range projects.Projects {
		id := print.Column{}
		name := print.Column{}

		id.Content = strconv.FormatInt(int64(project.ID), 10)
		id.FgColor = print.ID

		name.Content = project.Name
		name.ParentPad = true

		l.AddRow(project.ID, project.Parent.ID, id, name)
	}

	l.SetLimit(projects.Limit)
	l.SetOffset(projects.Offset)
	l.SetTotal(projects.TotalCount)

	editor.StartPage(r.Config.Pager, l.Render())
}

func cmdProjectList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Long:  "List all projects",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/projects.json?limit=1000&")
		},
	}

	// All
	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "List all projects",
		Long:  "List all projects",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/projects.json?")
		},
	})

	cmd.PersistentFlags().StringP(FLAG_QUERY, FLAG_QUERY_SHORT, "", "Query for projects with name")

	util.AddFlags(cmd)

	return cmd
}
