package issues

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

type issues struct {
	Issues []issue `json:"issues,omitempty"`
}

type issue struct {
	ID           int         `json:"id,omitempty"`
	Project      issueIdName `json:"project,omitempty"`
	Tracker      issueIdName `json:"tracker,omitempty"`
	Status       issueStatus `json:"status,omitempty"`
	Priority     issueIdName `json:"priority,omitempty"`
	Author       issueIdName `json:"author,omitempty"`
	AssignedTo   issueIdName `json:"assigned_to,omitempty"`
	FixedVersion issueIdName `json:"fixed_version,omitempty"`
	Subject      string      `json:"subject,omitempty"`
	Description  string      `json:"description,omitempty"`
	StartDate    string      `json:"start_date,omitempty"`
	DueDate      string      `json:"due_date,omitempty"`
	DoneRatio    int         `json:"done_ratio,omitempty"`
	IsPrivate    bool        `json:"is_private,omitempty"`
	//	EstimatedHours
	//	TotalEstimatedHours
	SpentHours      int    `json:"spent_hours,omitempty"`
	TotalSpentHours int    `json:"total_spent_hours,omitempty"`
	CreatedOn       string `json:"created_on,omitempty"`
	UpdatedOn       string `json:"updated_on,omitempty"`
	//	ClosedOn
}

type issueIdName struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type issueStatus struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}

func runList(r *config.Red_t) {
	var err error
	var body []byte
	issues := issues{}

	if body, err = api.ClientGET(r, "/issues.json"); err != nil {
		fmt.Println("Could not get response from client", err)
		return
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		fmt.Println("Could not parse and read response from server")
		return
	}

	for _, issue := range issues.Issues {
		println(issue.ID, issue.Subject)
	}
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
