package issues

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

/**
	"id":354,
	"project":{"id":79,"name":"PWDL (Personal Wireless Device Locator)"},
	"tracker":{"id":7,"name":"Task"},
	"status":{"id":2,"name":"In Progress","is_closed":false},
	"priority":{"id":3,"name":"High"},
	"author":{"id":5,"name":"Jeff Hägerman"},
	"assigned_to":{"id":5,"name":"Jeff Hägerman"},
	"fixed_version":{"id":90,"name":"v1.0"},
	"subject":"ECO-ESP-LIB struktur och bas",
	"description":"Strukturera ett basiskt behov av funktioner och komponenter för kunna lägga enkelt på en PWDL exempelvis eller dylikt i framtiden",
	"start_date":"2022-06-03",
	"due_date":null,
	"done_ratio":20,
	"is_private":false,
	"estimated_hours":null,
	"total_estimated_hours":null,
	"spent_hours":0.0,
	"total_spent_hours":0.0,
	"created_on":"2022-06-03T15:43:14Z",
	"updated_on":"2022-06-03T15:55:20Z",
	"closed_on":null
}
*/

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
