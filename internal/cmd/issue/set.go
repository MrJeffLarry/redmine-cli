package issue

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

func setIssueStatus(r *config.Red_t, id string) {
	var body []byte
	var status int
	var err error
	var statusHolder issueStatusHolder
	var issue newIssueHolder

	body, status, err = api.ClientGET(r, "/issue_statuses.json")
	print.Debug(r, status, string(body))
	if err != nil || status != 200 {
		print.Error("Could not get statuses from server, abort")
		return
	}

	if err := json.Unmarshal(body, &statusHolder); err != nil {
		print.Debug(r, status, err.Error())
		print.Error("Could not parse and read response from server")
	}

	fmt.Printf("Do check what status is allowed to change from status to status, as there is rules we can not read, so double check status got set\n\nSet new status for issue %s\n\n", text.FgGreen.Sprint(id))

	var idNames []util.IdName
	for _, status := range statusHolder.IssueStatus {
		idname := util.IdName{
			ID:   status.ID,
			Name: status.Name,
		}
		idNames = append(idNames, idname)
	}

	issue.Issue.StatusID, _ = terminal.Choose("Choose Status", idNames)

	/*
		fmt.Printf("Choose Status\n")
		for _, s := range statusHolder.IssueStatus {
			close := ""
			if s.IsClosed {
				close = "(Close)"
			}
			fmt.Printf("-> %s %s\n", s.Name, close)
		}

		for issue.Issue.StatusID == 0 {
			newStatus := terminal.WriteLine("New Status")
			for _, s := range statusHolder.IssueStatus {
				if s.Name == newStatus {
					issue.Issue.StatusID = s.ID
					break
				}
			}
			if issue.Issue.StatusID == 0 {
				print.Error("Status %s does not exist", newStatus)
			}
		}
	*/

	if body, err = json.Marshal(&issue); err != nil {
		print.Debug(r, 0, err.Error())
		print.Error("Could not compile request and send it, abort..")
		return
	}

	print.Debug(r, 0, string(body))

	if !terminal.Confirm("Set new status?") {
		return
	}

	body, status, err = api.ClientPUT(r, "/issues/"+id+".json", body)
	print.Debug(r, status, string(body))
	if err != nil || status != 204 {
		print.Error("%d Could not set new status..", status)
		return
	}

	print.OK("Status on issue #%s updated!", id)
}

func cmdIssueSetStatus(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "set status",
		Long:    "set status on an issue",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			id := cmd.Flags().Arg(0)

			if len(id) <= 0 {
				fmt.Println("Please specify what issue you would like to set, usage: set [id]")
				return
			}

			setIssueStatus(r, id)
		},
	}
	return cmd
}

func cmdIssueSet(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [id]",
		Short:   "set issue",
		Long:    "set an issue",
		Aliases: []string{"s"},
	}

	cmd.AddCommand(cmdIssueSetStatus(r))

	return cmd
}
