package issue

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
	HEAD_ID             = "ID"
	HEAD_STATUS         = "STATUS"
	HEAD_TRACKER        = "TRACKER"
	HEAD_TARGET_VERSION = "TARGET VERSION"
	HEAD_PRIORITY       = "PRIORITY"
	HEAD_PROJECT        = "PROJECT"
	HEAD_SUBJECT        = "SUBJECT"

	FLAG_DISPLAY_PROJECT = "project"
	FLAG_QUERY           = "query"
	FLAG_QUERY_SHORT     = "q"
)

func displayListGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int
	var dispProject bool

	head := []string{HEAD_ID, HEAD_TRACKER, HEAD_STATUS, HEAD_PRIORITY, HEAD_TARGET_VERSION, HEAD_SUBJECT}
	sort := []string{"id", "status", "priority", "subject"}
	issues := issues{}
	projectID := 0

	if !r.All {
		projectID = r.Config.ProjectID
	}

	if dispProject, _ = cmd.Flags().GetBool(FLAG_DISPLAY_PROJECT); dispProject {
		head = append(head, HEAD_PROJECT)
		sort = append(sort, "project")
	}

	if query, _ := cmd.Flags().GetString(FLAG_QUERY); query != "" {
		path += "subject=~" + query + "&"
	}

	path += util.ParseFlags(cmd, projectID, sort)

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
		tracker := print.Column{}
		status := print.Column{}
		priority := print.Column{}
		subject := print.Column{}
		project := print.Column{}
		targetVersion := print.Column{}

		id.Content = strconv.FormatInt(int64(issue.ID), 10)
		id.FgColor = print.ID
		id.ParentPad = true

		tracker.Content = issue.Tracker.Name
		status.Content = issue.Status.Name
		priority.Content = issue.Priority.Name
		subject.Content = issue.Subject
		project.Content = issue.Project.Name
		targetVersion.Content = issue.FixedVersion.Name

		if dispProject {
			l.AddRow(issue.ID, issue.Parent.ID, id, tracker, status, priority, targetVersion, subject, project)
		} else {
			l.AddRow(issue.ID, issue.Parent.ID, id, tracker, status, priority, targetVersion, subject)
		}
	}
	l.SetLimit(issues.Limit)
	l.SetOffset(issues.Offset)
	l.SetTotal(issues.TotalCount)

	editor.StartPage(r.Config.Pager, l.Render())
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
			r.Config.ProjectID = 0 // ignore ID if we want too see all
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

	cmd.PersistentFlags().Bool(FLAG_DISPLAY_PROJECT, false, "Display project column")
	cmd.PersistentFlags().StringP(FLAG_QUERY, FLAG_QUERY_SHORT, "", "Query for issues with subject")

	util.AddFlags(cmd)

	return cmd
}
