package issue

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/editor"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

func writeLine(pre string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(pre, ": ")
	text, _ := reader.ReadString('\n')
	return text
}

func displayCreateIssue(r *config.Red_t, cmd *cobra.Command, path string) {
	var projectIde string
	hold := true
	issue := issue{}

	if len(r.RedmineProject) > 0 {
		projectIde = r.RedmineProject
	}

	if proIde, _ := cmd.Flags().GetString("project"); len(proIde) > 0 {
		projectIde = proIde
	}

	if len(projectIde) == 0 {
		fmt.Println("Project identity is missing, please use `--project project-identity` or use local override .red/config.json, or global project")
		return
	}

	fmt.Printf("Create new issue in project %s\n\n", text.FgGreen.Sprint(projectIde))

	issue.Subject = writeLine("Subject")

	for hold {
		writeBody := writeLine("Write body? (y/n)")
		if strings.Contains(writeBody, "y") {
			issue.Description = editor.StartEdit("")
			hold = false
		} else if strings.Contains(writeBody, "n") {
			hold = false
		} else {
			print.Error("%s: %s", "No valid input, valid (y=yes or n=no)", writeBody)
		}
	}

	fmt.Printf("Subject %s\nDescription: %s\n", issue.Subject, issue.Description)
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

	cmd.PersistentFlags().StringP("project", "p", "", "What project identity should be used for the new issue")

	return cmd
}
