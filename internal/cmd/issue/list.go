package issue

import (
	"encoding/json"
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/spf13/cobra"
)

const (
	HEAD_ID       = "ID"
	HEAD_STATUS   = "STATUS"
	HEAD_PRIORITY = "PRIORITY"
	HEAD_PROJECT  = "PROJECT"
	HEAD_SUBJECT  = "SUBJECT"
)

func displayListGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int
	head := []string{HEAD_ID, HEAD_STATUS, HEAD_PRIORITY, HEAD_PROJECT, HEAD_SUBJECT}

	issues := issues{}
	projectID := 0

	if !r.All {
		projectID = r.RedmineProjectID
	}

	path += util.ParseFlags(cmd, projectID, []string{"id", "status", "priority", "project", "subject"})

	print.Debug(r, path)

	if body, status, err = api.ClientGET(r, path); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.Debug(r, "%d %s", status, string(body))

	if err := json.Unmarshal(body, &issues); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	l := print.NewList(head...)

	for _, issue := range issues.Issues {
		id := print.Column{}
		status := print.Column{}
		priority := print.Column{}
		project := print.Column{}
		subject := print.Column{}

		id.Content = strconv.FormatInt(issue.ID, 10)
		id.FgColor = print.ID

		status.Content = issue.Status.Name
		priority.Content = issue.Priority.Name
		project.Content = issue.Project.Name
		subject.Content = issue.Subject

		l.AddRow(id, status, priority, project, subject)
	}
	l.SetLimit(issues.Limit)
	l.SetOffset(issues.Offset)
	l.SetTotal(issues.TotalCount)

	l.Render()
}

func cmdIssueList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Long:  "List all issues",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/issues.json?")
		},
	}

	// All
	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "List all issues",
		Long:  "List all issues and ignores project ID",
		Run: func(cmd *cobra.Command, args []string) {
			r.RedmineProjectID = 0 // ignore ID if we want too see all
			displayListGET(r, cmd, "/issues.json?")
		},
	})

	// Me
	cmd.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "List all my issues",
		Long:  "List all my issues",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/issues.json?assigned_to_id=me&")
		},
	})

	util.AddFlags(cmd)

	return cmd
}
