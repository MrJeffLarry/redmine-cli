package issue

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/project"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

func displayCreateIssue(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var projectID int
	var trackers []util.IdName

	issue := newIssueHolder{}

	if r.RedmineProjectID > 0 {
		projectID = r.RedmineProjectID
	}

	if proID, _ := cmd.Flags().GetInt("project"); proID > 0 {
		projectID = proID
	}

	if projectID <= 0 {
		fmt.Println("Project ID is missing, please use `--project 20` or use local override .red/config.json, or global project")
		return
	}

	issue.Issue.ProjectID = int64(projectID)

	if trackers, err = project.GetTrackers(r, projectID); err != nil {
		print.Error(err.Error())
		return
	}

	fmt.Printf("Create new issue in project %s\n\n", text.FgGreen.Sprint(projectID))

	issue.Issue.TrackerID, _ = terminal.Choose("Tracker", trackers)
	issue.Issue.Subject = terminal.WriteLineReq("Subject", 1)
	if terminal.Confirm("Write Body?") {
		issue.Issue.Description = editor.StartEdit("")
	}

	body, err := json.Marshal(issue)
	if err != nil {
		print.Debug(r, 0, err.Error())
		print.Error("Could not compose issue..")
		return
	}

	print.Debug(r, 0, string(body))

	if !terminal.Confirm("Create issue?") {
		return
	}

	if body, status, err := api.ClientPOST(r, "/issues.json", body); err != nil || status != 201 {
		print.Debug(r, 0, string(body))
		print.Error("%d Could not send issue", status)
		return
	}
	print.OK("Issue created!")
}

func cmdIssueCreate(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create issue",
		Long:    "Create an issue",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			displayCreateIssue(r, cmd, "/issues.json")
		},
	}

	cmd.PersistentFlags().IntP("project", "p", -1, "What project ID should be used for the new issue")

	return cmd
}
