package issue

import (
	"encoding/json"
	"errors"
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

func cmdIssueEditIssueAssign(r *config.Red_t, projectID int, issue *newIssueHolder) error {
	var err error
	var idNames []util.IdName

	if idNames, err = project.GetAssigns(r, projectID); err != nil {
		print.Error(err.Error())
	}

	id, _ := terminal.Choose("Assign", idNames)

	if id >= 0 {
		issue.Issue.AssignedToID = id
	}
	return nil
}

func cmdIssueEditIssueNote(r *config.Red_t, issue *newIssueHolder) error {
	issue.Issue.Notes = editor.StartEdit("")
	return nil
}

func cmdIssueEditIssueStatus(r *config.Red_t, issue *newIssueHolder) error {
	var body []byte
	var status int
	var err error
	var statusHolder issueStatusHolder

	body, status, err = api.ClientGET(r, "/issue_statuses.json")
	print.Debug(r, "%d %s", status, string(body))
	if err != nil || status != 200 {
		return errors.New("Could not get statuses from server, abort")
	}

	if err := json.Unmarshal(body, &statusHolder); err != nil {
		print.Debug(r, err.Error())
		return errors.New("Could not parse and read response from server")
	}

	fmt.Println(text.FgHiBlack.Sprint("Do check what status is allowed to change from status to status, as there is rules we can not read, so double check status got set"))

	var idNames []util.IdName
	for _, status := range statusHolder.IssueStatus {
		idname := util.IdName{
			ID:   status.ID,
			Name: status.Name,
		}
		idNames = append(idNames, idname)
	}

	issue.Issue.StatusID, _ = terminal.Choose("Choose Status", idNames)
	return nil
}

func cmdIssueEditIssuePriority(r *config.Red_t, issue *newIssueHolder) error {
	var body []byte
	var status int
	var err error
	var priorityHolder issuePrioritiesHolder

	body, status, err = api.ClientGET(r, "/enumerations/issue_priorities.json")
	print.Debug(r, "%d %s", status, string(body))
	if err != nil || status != 200 {
		return errors.New("Could not get statuses from server, abort")
	}

	if err := json.Unmarshal(body, &priorityHolder); err != nil {
		print.Debug(r, err.Error())
		return errors.New("Could not parse and read response from server")
	}

	var idNames []util.IdName
	for _, prio := range priorityHolder.IssuePriorities {
		idname := util.IdName{
			ID:   prio.ID,
			Name: prio.Name,
		}
		idNames = append(idNames, idname)
	}

	issue.Issue.PriorityID, _ = terminal.Choose("Choose Priority", idNames)
	return nil
}

func cmdIssueEditIssueTracker(r *config.Red_t, projectID int, issue *newIssueHolder) error {
	var err error
	var trackers []util.IdName

	if trackers, err = project.GetTrackers(r, projectID); err != nil {
		return err
	}

	issue.Issue.TrackerID, _ = terminal.Choose("Tracker", trackers)
	return nil
}

func cmdIssueEditIssueDescription(r *config.Red_t, description string, issue *newIssueHolder) error {
	issue.Issue.Description = editor.StartEdit(description)
	return nil
}

func cmdIssueEditIssueSubject(r *config.Red_t, subject string, issue *newIssueHolder) error {
	var err error
	if issue.Issue.Subject, err = terminal.PromptString("Subject", subject); err != nil {
		return err
	}
	return nil
}

func cmdIssueEditIssueSave(r *config.Red_t, path string, issue *newIssueHolder) {
	var body []byte
	var err error
	var errList util.Errors

	body, err = json.Marshal(issue)
	if err != nil {
		print.Debug(r, err.Error())
		print.Error("Could not compose issue..")
		return
	}

	print.Debug(r, string(body))

	if !terminal.Confirm("Confirm save issue?") {
		return
	}

	if body, status, err := api.ClientPUT(r, path, body); err != nil || status != 204 {
		if err = json.Unmarshal(body, &errList); err != nil {
			print.Error("%d Could not read error response from server: %s", status, string(body))
			return
		}
		print.Debug(r, string(body))
		print.Error("%d Could not save issue: %v", status, errList.Errors)
		return
	}
	print.OK("Issue saved!")
}

func cmdIssueEditIssueDebug(r *config.Red_t, issue *newIssueHolder) {
	print.Debug(r, "\n"+
		"Subject: %s\n"+
		"TrackerID: %d\n"+
		"StatusID: %d\n"+
		"Notes: %s\n"+
		"Description: %s\n",
		issue.Issue.Subject,
		issue.Issue.TrackerID,
		issue.Issue.StatusID,
		issue.Issue.Notes,
		issue.Issue.Description,
	)
}

func cmdIssueEditIssue(r *config.Red_t, cmd *cobra.Command, id, path string) {
	var err error
	var body []byte
	var status int
	var viewIssue viewIssue
	chooses := []string{
		FIELD_SUBJECT,
		FIELD_DESCRIPTION,
		FIELD_STATUS,
		FIELD_PRIORITY,
		FIELD_TRACKER,
		FIELD_NOTE,
		FIELD_ASSIGN,
		FIELD_PREVIEW,
		FIELD_SAVE,
		FIELD_EXIT}

	if r.Debug {
		chooses = append(chooses, FIELD_DEBUG)
	}

	issue := newIssueHolder{}

	if body, status, err = api.ClientGET(r, path); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.Debug(r, "%d %s", status, string(body))

	if err := api.StatusCode(status); err != nil {
		print.Error(err.Error())
		return
	}

	if err := json.Unmarshal(body, &viewIssue); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	fmt.Printf("Edit issue #%s\n\n", text.FgGreen.Sprint(id))

	for {
		choose, i := terminal.ChooseString("Issue", chooses)
		if i == -1 {
			if !terminal.Confirm("Exit") {
				continue
			}
			return
		}
		switch choose {
		case FIELD_SUBJECT:
			if err = cmdIssueEditIssueSubject(r, viewIssue.Issue.Subject, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_STATUS:
			if err = cmdIssueEditIssueStatus(r, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_PRIORITY:
			if err = cmdIssueEditIssuePriority(r, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_TRACKER:
			if err = cmdIssueEditIssueTracker(r, viewIssue.Issue.Project.ID, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_DESCRIPTION:
			if err = cmdIssueEditIssueDescription(r, viewIssue.Issue.Description, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_NOTE:
			if err = cmdIssueEditIssueNote(r, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_ASSIGN:
			if err = cmdIssueEditIssueAssign(r, viewIssue.Issue.Project.ID, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_SAVE:
			cmdIssueEditIssueSave(r, path, &issue)
			return
		case FIELD_DEBUG:
		case FIELD_PREVIEW:
			cmdIssueEditIssueDebug(r, &issue)
		case FIELD_EXIT:
			if !terminal.Confirm("Exit") {
				continue
			}
			return
		}
	}
}

func cmdIssueEdit(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit [id]",
		Short:   "edit issue",
		Long:    "edit an issue",
		Aliases: []string{"e"},
		Run: func(cmd *cobra.Command, args []string) {
			id := cmd.Flags().Arg(0)

			if !util.CheckID(id) {
				fmt.Println("Please specify what issue you would like to edit, usage: edit [id]")
				return
			}
			cmdIssueEditIssue(r, cmd, id, "/issues/"+id+".json")
		},
	}

	return cmd
}
