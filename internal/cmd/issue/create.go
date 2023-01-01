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

func cmdIssueCreateSave(r *config.Red_t, issue *newIssueHolder) bool {
	body, err := json.Marshal(issue)
	if err != nil {
		print.Debug(r, err.Error())
		print.Error("Could not compose issue..")
		return false
	}

	print.Debug(r, string(body))

	response, status, err := api.ClientPOST(r, "/issues.json", body)
	if err != nil {
		print.Error("Error %v", err)
		return false
	}

	print.Debug(r, string(response))

	if status != 201 {
		errors := api.ParseResponseError(response)
		print.Error("Could not save issue")
		print.Error("%s", errors.Errors)
		return false
	}

	print.OK("Issue created!")
	return true
}

func displayCreateIssue(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var projectID int
	var trackers []util.IdName
	var versions []util.IdName
	chooses := []string{
		FIELD_SAVE,
		FIELD_VERSION,
		FIELD_PARENT_ID,
		FIELD_EXIT}

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

	issue.Issue.ProjectID = projectID

	if trackers, err = project.GetTrackers(r, projectID); err != nil {
		print.Error(err.Error())
		return
	}

	fmt.Printf("Create new issue in project %s\n\n", text.FgGreen.Sprint(projectID))

	issue.Issue.TrackerID, _ = terminal.Choose("Tracker", trackers)
	issue.Issue.Subject, _ = terminal.PromptStringRequire("Subject", "")
	if terminal.Confirm("Write Body") {
		issue.Issue.Description = editor.StartEdit("")
	}

	//
	for {
		choose, i := terminal.ChooseString("Options", chooses)
		if i == -1 {
			if !terminal.Confirm("Exit") {
				continue
			}
			return
		}

		switch choose {
		case FIELD_SAVE:
			if cmdIssueCreateSave(r, &issue) {
				return
			}
		case FIELD_VERSION:
			if versions, err = project.GetVersions(r, projectID); err != nil {
				print.Error(err.Error())
			}
			issue.Issue.FixedVersionID, _ = terminal.Choose("Version", versions)
		case FIELD_PARENT_ID:
			parentID, _ := terminal.PromptInt("Parent ID (-1 means none)", -1)
			if parentID > 0 {
				issue.Issue.ParentIssueID = parentID
			}
		case FIELD_ASSIGN:
			if r.RedmineUserID > 0 && terminal.Confirm("Assign to me") {
				issue.Issue.AssignedToID = r.RedmineUserID
			}
		case FIELD_EXIT:
			if !terminal.Confirm("Exit") {
				continue
			}
			return
		}
	}
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
