package issue

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func writeLine(pre string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(pre, ": ")
	text, _ := reader.ReadString('\n')
	return text
}

func displayCreateIssue(r *config.Red_t, cmd *cobra.Command, path string) {
	var projectID int16
	var err error
	issue := issue{}

	if projectID, err = cmd.Flags().GetInt16("project"); err != nil || projectID < 0 {
		fmt.Println("Project id is missing, please use `--project 2` for project 2")
		return
	}

	fmt.Print("Create new issue\n\n")

	issue.Subject = writeLine("Subject")
	issue.Description = writeLine("Description")

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

	cmd.PersistentFlags().Int16P("project", "p", -1, "What project id should the new issue use")

	return cmd
}
