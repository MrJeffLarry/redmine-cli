package issues

import (
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

type issue struct {
	ID           int
	Project      issueIdName
	Tracker      issueIdName
	Status       issueStatus
	Priority     issueIdName
	Author       issueIdName
	AssignedTo   issueIdName
	FixedVersion issueIdName
}

type issueIdName struct {
	ID   int
	Name string
}

type issueStatus struct {
	ID       int
	Name     string
	IsClosed bool
}

func runList(r *config.Red_t) {
	body, err := api.ClientGET(r, "/issues.json")
	fmt.Println("Running list", string(body), err)
}

func cmdIssuesList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Long:  "List all issues",
		Run: func(cmd *cobra.Command, args []string) {
			runList(r)
		},
	}
	return cmd
}
