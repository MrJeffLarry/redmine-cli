package issues

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

type issues struct {
	Issues     []issue `json:"issues,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type issue struct {
	ID                  int64       `json:"id,omitempty"`
	Project             issueIdName `json:"project,omitempty"`
	Tracker             issueIdName `json:"tracker,omitempty"`
	Status              issueStatus `json:"status,omitempty"`
	Priority            issueIdName `json:"priority,omitempty"`
	Author              issueIdName `json:"author,omitempty"`
	AssignedTo          issueIdName `json:"assigned_to,omitempty"`
	FixedVersion        issueIdName `json:"fixed_version,omitempty"`
	Subject             string      `json:"subject,omitempty"`
	Description         string      `json:"description,omitempty"`
	StartDate           string      `json:"start_date,omitempty"`
	DueDate             string      `json:"due_date,omitempty"`
	DoneRatio           int         `json:"done_ratio,omitempty"`
	IsPrivate           bool        `json:"is_private,omitempty"`
	EstimatedHours      float32     `json:"estimated_hours,omitempty"`
	TotalEstimatedHours float32     `json:"total_estimated_hours,omitempty"`
	SpentHours          float32     `json:"spent_hours,omitempty"`
	TotalSpentHours     float32     `json:"total_spent_hours,omitempty"`
	CreatedOn           string      `json:"created_on,omitempty"`
	UpdatedOn           string      `json:"updated_on,omitempty"`
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

const (
	FLAG_ORDER     = "order"
	FLAG_ORDER_ASC = "asc"
	FLAG_ORDER_DES = "des"

	FLAG_LIMIT = "limit"
)

func countDigi(i int64) (count int) {
	for i > 0 {
		i = i / 10
		count++
	}
	return
}

func parseFlags(cmd *cobra.Command, path string) string {
	//	order, err := cmd.Flags().GetString(FLAG_ORDER)
	return path
}

func displayListGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int
	var statusLen int
	var idLen int
	//	var subjectLen int
	issues := issues{}

	path = parseFlags(cmd, path)

	if body, status, err = api.ClientGET(r, path); err != nil {
		fmt.Println(status, "Could not get response from client", err)
		return
	}

	print.PrintDebug(r, status, string(body))

	if err := json.Unmarshal(body, &issues); err != nil {
		fmt.Println(err)
		fmt.Println(status, "Could not parse and read response from server")
		return
	}

	for _, issue := range issues.Issues {
		sLen := len(issue.Status.Name)
		iLen := countDigi(issue.ID)

		if sLen > statusLen {
			statusLen = sLen
		}

		if iLen > idLen {
			idLen = iLen
		}
	}

	for _, issue := range issues.Issues {
		sLeft := statusLen - len(issue.Status.Name)
		iLeft := idLen - countDigi(issue.ID)

		status := issue.Status.Name + strings.Repeat(" ", sLeft) // strconv.Itoa(sLeft)
		idPad := strings.Repeat(" ", iLeft)

		fmt.Printf("#%d%s - %s - %s\n", issue.ID, idPad, status, issue.Subject)
	}
}

func cmdIssuesList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Long:  "List all issues",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/issues.json?")
		},
	}

	// All
	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "List all issues",
		Long:  "List all issues",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/issues.json?")
		},
	})

	// Me
	cmd.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "List all my issues",
		Long:  "List all my issues",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/issues.json?assigned_to_id=me")
		},
	})

	cmd.PersistentFlags().String("order", "", "Order on id_ASC or id_DES")
	cmd.PersistentFlags().String("sort", "", "")

	return cmd
}
