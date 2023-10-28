package issue

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

const (
	FLAG_JOURNALS   = "journals"
	FLAG_JOURNALS_S = "j"
)

func displayProgressBar(ratio int) string {
	done := (ratio / 10)
	progress := "["
	progress += strings.Repeat("=", done)
	progress += strings.Repeat(" ", 10-done)
	progress += "]"
	return progress
}

func displayIssue(r *config.Red_t, i issue, journalFlag bool) {
	closed := text.FgGreen.Sprint("OPEN")

	if i.Status.IsClosed {
		closed = text.FgRed.Sprint("CLOSED")
	}

	sid := strconv.FormatInt(int64(i.ID), 10)

	response := fmt.Sprintf(
		"------------ %s %s - %s [%s] ---------\n"+
			text.FgGreen.Sprint("Author")+" %s\n"+
			text.FgGreen.Sprint("Start Date")+" %s\n"+
			text.FgGreen.Sprint("Due Date")+" %s\n"+
			text.FgGreen.Sprint("Done")+" %s %d%%\n\n"+
			text.FgGreen.Sprint("Assigned")+" %s\n"+
			text.FgGreen.Sprint("Created")+" %s\n"+
			text.FgGreen.Sprint("Project")+" %s\n"+
			text.FgGreen.Sprint("Version")+" %s\n"+
			text.FgGreen.Sprint("Status")+" %s\n"+
			text.FgGreen.Sprint("Priority")+" %s\n"+
			text.FgGreen.Sprint("Description")+"\n"+
			"\n%s\n\n",
		text.FgYellow.Sprint(i.Tracker.Name),
		print.PrintID(i.ID),
		i.Subject,
		closed,
		i.Author.Name,
		i.StartDate,
		i.DueDate,
		displayProgressBar(i.DoneRatio),
		i.DoneRatio,
		i.AssignedTo.Name,
		i.CreatedOn,
		i.Project.Name,
		i.FixedVersion.Name,
		i.Status.Name,
		i.Priority.Name,
		i.Description,
	)

	if journalFlag {
		for _, journal := range i.Journals {
			status := ""
			notes := ""
			for _, detail := range journal.Details {
				status += text.FgGreen.Sprint("Update") + " "
				status += detail.Name + " changed from "
				status += detail.OldValue + " to "
				status += detail.NewValue + "\n"
			}
			if len(journal.Notes) > 0 {
				notes = text.FgGreen.Sprint("Notes") + " "
				notes += journal.Notes + "\n"
			}

			response += fmt.Sprintf(" %s %s\n"+
				"%s"+
				"%s\n",
				text.FgGreen.Sprintf("#%d", journal.ID),
				journal.User.Name,
				text.FgHiBlack.Sprint(notes),
				text.FgHiBlack.Sprint(status),
			)
		}
	}

	response += fmt.Sprintln(text.FgHiBlack.Sprintf("View issue: %s", r.Config.Server+"/issues/"+sid))

	editor.StartPage(r.Config.Pager, response)
}

func displayIssueGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int

	viewIssue := viewIssue{}

	if body, status, err = api.ClientGET(r, path); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.Debug(r, "%d %s", status, string(body))

	if err := api.StatusCode(status); err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(body, &viewIssue); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	journals, _ := cmd.Flags().GetBool(FLAG_JOURNALS)

	displayIssue(r, viewIssue.Issue, journals)
}

func cmdIssueView(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [id]",
		Short: "View issue",
		Long:  "View issue information and details",
		Run: func(cmd *cobra.Command, args []string) {
			id := cmd.Flags().Arg(0)
			include := "?include=allowed_statuses"

			if len(id) <= 0 {
				fmt.Println("Please specify what issue you would like to view, usage: view [id]")
				return
			}

			if journals, _ := cmd.Flags().GetBool(FLAG_JOURNALS); journals {
				include += ",journals"
			}

			displayIssueGET(r, cmd, "/issues/"+id+".json"+include)
		},
	}

	cmd.PersistentFlags().BoolP(FLAG_JOURNALS, FLAG_JOURNALS_S, false, "Display journals")

	return cmd
}
