package issue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/global"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/project"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

func cmdIssueEditIssueTargetVersion(r *config.Red_t, projectID int) (util.IdName, error) {
	var err error
	var idName util.IdName
	var idNames []util.IdName

	if idNames, err = project.GetVersions(r, projectID); err != nil {
		return idName, err
	}

	idName.ID, idName.Name = r.Term.Choose("Target version", idNames)

	if idName.ID < 0 {
		return idName, errors.New("target version ID not valid")
	}

	return idName, nil
}

func cmdIssueEditIssueAssign(r *config.Red_t, projectID int) (util.IdName, error) {
	var err error
	var idName util.IdName
	var idNames []util.IdName

	if idNames, err = project.GetAssigns(r, projectID); err != nil {
		return idName, err
	}

	idName.ID, idName.Name = r.Term.Choose("Assign", idNames)

	if idName.ID < 0 {
		return idName, errors.New("assigne ID not valid")
	}

	return idName, nil
}

func cmdIssueEditIssueNote(r *config.Red_t, issue *newIssueHolder) error {
	issue.Issue.Notes = editor.StartEdit(r.Config.Editor, "")
	return nil
}

func cmdIssueEditIssueStatus(r *config.Red_t, allowedStatus []global.IssueStatus) (global.IssueStatus, error) {
	var err error
	var idNames []util.IdName
	var issueStatus global.IssueStatus

	if len(allowedStatus) > 0 {
		for _, item := range allowedStatus {
			row := util.IdName{
				ID:   item.ID,
				Name: item.Name,
			}
			idNames = append(idNames, row)
		}
	} else {
		fmt.Println(text.FgHiBlack.Sprint("allowed_statuses is missing or empty, Redmine 5.0.x > supports allowed_statuses so we can read what status possible to change to, we will serve you all the global statuses instead"))

		if idNames, err = global.GetIssueStatus(r); err != nil {
			return issueStatus, err
		}
	}

	issueStatus.ID, issueStatus.Name = r.Term.Choose("Choose Status", idNames)
	return issueStatus, nil
}

func cmdIssueEditIssuePriority(r *config.Red_t) (util.IdName, error) {
	var err error
	var idName util.IdName
	var idNames []util.IdName

	if idNames, err = global.GetPriorities(r); err != nil {
		return idName, err
	}

	idName.ID, idName.Name = r.Term.Choose("Choose Priority", idNames)
	return idName, nil
}

func cmdIssueEditIssueTracker(r *config.Red_t, projectID int) (util.IdName, error) {
	var err error
	var idName util.IdName
	var trackers []util.IdName

	if trackers, err = project.GetTrackers(r, projectID); err != nil {
		return idName, err
	}

	idName.ID, idName.Name = r.Term.Choose("Tracker", trackers)
	return idName, nil
}

func cmdIssueEditIssueSave(r *config.Red_t, id, path string, issue *newIssueHolder) bool {
	var body []byte
	var err error
	var errList util.Errors

	body, err = json.Marshal(issue)
	if err != nil {
		print.Debug(r, err.Error())
		print.Error("Could not compose issue..")
		return false
	}

	print.Debug(r, string(body))

	if body, status, err := api.ClientPUT(r, path, body); err != nil || status != 204 {
		if err = json.Unmarshal(body, &errList); err != nil {
			print.Error("%d Could not read error response from server: %s", status, string(body))
			return false
		}

		print.Debug(r, string(body))
		print.Error("%d Could not save issue: %v", status, errList.Errors)
		return false
	}

	print.OK("Issue #%s saved!", id)

	return true
}

func cmdIssueEditIssue(r *config.Red_t, cmd *cobra.Command, id, path string) {
	var err error
	var body []byte
	var status int
	var viewIssue viewIssue

	chooses := []string{
		FIELD_SAVE,
		FIELD_SUBJECT,
		FIELD_DESCRIPTION,
		FIELD_STATUS,
		FIELD_PRIORITY,
		FIELD_TRACKER,
		FIELD_NOTE,
		FIELD_ASSIGN,
		FIELD_TARGET_VERSION,
		FIELD_PREVIEW,
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

	fmt.Printf("Edit issue %s - %s\n\n", print.PrintID(viewIssue.Issue.ID), viewIssue.Issue.Subject)

	for {
		choose, i := r.Term.ChooseString("Issue", chooses)
		if i == -1 {
			if !r.Term.Confirm("Exit") {
				continue
			}
			return
		}
		switch choose {
		case FIELD_SUBJECT:
			subject, err := r.Term.PromptString("Subject", viewIssue.Issue.Subject)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.Subject = subject
				viewIssue.Issue.Subject = subject
			}
		case FIELD_STATUS:
			idName, err := cmdIssueEditIssueStatus(r, viewIssue.Issue.AllowedStatuses)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.StatusID = idName.ID
				viewIssue.Issue.Status = idName
			}
		case FIELD_PRIORITY:
			idName, err := cmdIssueEditIssuePriority(r)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.PriorityID = idName.ID
				viewIssue.Issue.Priority = idName
			}
		case FIELD_TRACKER:
			idName, err := cmdIssueEditIssueTracker(r, viewIssue.Issue.Project.ID)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.TrackerID = idName.ID
				viewIssue.Issue.Tracker = idName
			}
		case FIELD_DESCRIPTION:
			issue.Issue.Description = editor.StartEdit(r.Config.Editor, viewIssue.Issue.Description)
			viewIssue.Issue.Description = issue.Issue.Description
		case FIELD_NOTE:
			if err = cmdIssueEditIssueNote(r, &issue); err != nil {
				print.Error(err.Error())
			}
		case FIELD_ASSIGN:
			idName, err := cmdIssueEditIssueAssign(r, viewIssue.Issue.Project.ID)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.AssignedToID = idName.ID
				viewIssue.Issue.AssignedTo = idName
			}
		case FIELD_TARGET_VERSION:
			idName, err := cmdIssueEditIssueTargetVersion(r, viewIssue.Issue.Project.ID)

			if err != nil {
				print.Error(err.Error())
			} else {
				issue.Issue.FixedVersionID = idName.ID
				viewIssue.Issue.FixedVersion = idName
			}
		case FIELD_SAVE:
			if cmdIssueEditIssueSave(r, id, path, &issue) {
				return
			}
		case FIELD_DEBUG:
		case FIELD_PREVIEW:
			displayIssue(r, viewIssue.Issue, false)
		case FIELD_EXIT:
			if !r.Term.Confirm("Exit") {
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
			cmdIssueEditIssue(r, cmd, id, "/issues/"+id+".json?include=allowed_statuses")
		},
	}

	return cmd
}
