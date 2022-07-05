package issue

import (
	"encoding/json"
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/spf13/cobra"
)

func addNote(r *config.Red_t, path string, cmd *cobra.Command) {
	var err error
	var body []byte
	var status int
	issue := viewIssue{}

	issue.Issue.Notes = editor.StartEdit("")

	print.PrintDebug(r, 0, "Path:"+path)
	print.PrintDebug(r, 0, "Notes: "+issue.Issue.Notes)

	if !terminal.Confirm("Send Note?") {
		return
	}

	if body, err = json.Marshal(&issue); err != nil {
		return
	}

	print.PrintDebug(r, 0, string(body))

	if body, status, err = api.ClientPUT(r, path, body); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	} else {
		print.PrintDebug(r, status, string(body))
		print.OK("Note sent to issue!")
	}
}

func cmdIssueNoteAdd(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [id]",
		Short: "add note on issues",
		Long:  "add note on issues",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			if id, err = strconv.Atoi(cmd.Flags().Arg(0)); err != nil {
				print.Error("No valid issue id found")
				return
			}

			addNote(r, "/issues/"+strconv.Itoa(id)+".json", cmd)
		},
	}
	return cmd
}

func cmdIssueNote(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Note on issues",
		Long:  "Note on issues",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
			}
		},
	}

	cmd.AddCommand(cmdIssueNoteAdd(r))

	return cmd
}
