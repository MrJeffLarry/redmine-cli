package issues

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

type viewIssue struct {
	Issue issue `json:"issue,omitempty"`
}

func displayProgressBar(ratio int) string {
	done := (ratio / 10)
	progress := "["
	progress += strings.Repeat("=", done)
	progress += strings.Repeat(" ", 10-done)
	progress += "]"
	return progress
}

func displayIssueGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int

	viewIssue := viewIssue{}

	if body, status, err = api.ClientGET(r, path); err != nil {
		fmt.Println(status, "Could not get response from client", err)
		return
	}

	print.PrintDebug(r, status, string(body))

	if err := api.StatusCode(status); err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(body, &viewIssue); err != nil {
		fmt.Println(body)
		fmt.Println(err)
		fmt.Println(status, "Could not parse and read response from server")
		return
	}

	issue := viewIssue.Issue
	closed := "OPEN"

	if issue.Status.IsClosed {
		closed = "CLOSED"
	}

	sid := strconv.FormatInt(issue.ID, 10)

	fmt.Printf("\n"+
		"%s #%d %s [%s]\n"+
		"---------------------------------\n"+
		"Start Date: %s\n"+
		"Due Date: %s\n"+
		"Done: %s %d%%\n\n"+
		"Assigned: %s\n"+
		"Created: %s\n"+
		"Project: %s\n"+
		"Version: %s\n"+
		"Status: %s\n"+
		"Priority: %s\n"+
		"---------------------------------\n"+
		"\n%s\n\n"+
		"---------------------------------\n"+
		"View issue: %s\n",
		issue.Tracker.Name,
		issue.ID,
		issue.Subject,
		closed,
		issue.StartDate,
		issue.DueDate,
		displayProgressBar(issue.DoneRatio),
		issue.DoneRatio,
		issue.AssignedTo.Name,
		issue.CreatedOn,
		issue.Project.Name,
		issue.FixedVersion.Name,
		issue.Status.Name,
		issue.Priority.Name,
		issue.Description,
		r.RedmineURL+"/issues/"+sid,
	)
}

func cmdIssuesIssue(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [id]",
		Short: "View issue",
		Long:  "View issue information and details",
		Run: func(cmd *cobra.Command, args []string) {
			id := cmd.Flags().Arg(0)

			if len(id) <= 0 {
				fmt.Println("Please specify what issue you would like to view, usage: view [id]")
				return
			}

			displayIssueGET(r, cmd, "/issues/"+id+".json")
		},
	}

	return cmd
}
