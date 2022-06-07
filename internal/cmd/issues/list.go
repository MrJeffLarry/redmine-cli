package issues

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

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
	var subjectLen int
	var projectLen int

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
		subLen := len(issue.Subject)
		proLen := len(issue.Project.Name)
		sLen := len(issue.Status.Name)
		iLen := countDigi(issue.ID)

		if sLen > statusLen {
			statusLen = sLen
		}

		if iLen > idLen {
			idLen = iLen
		}

		if subLen > subjectLen {
			subjectLen = subLen
		}

		if proLen > projectLen {
			projectLen = proLen
		}
	}

	fmt.Printf("%s  %s  %s  %s\n",
		"ID"+strings.Repeat(" ", int(math.Abs(float64(idLen-len("ID"))))),
		"STATUS"+strings.Repeat(" ", int(math.Abs(float64(statusLen-len("STATUS"))))),
		"PROJECT"+strings.Repeat(" ", int(math.Abs(float64(projectLen-len("PROJECT"))))),
		"SUBJECT",
	)

	for _, issue := range issues.Issues {
		sLeft := statusLen - len(issue.Status.Name)
		iLeft := idLen - countDigi(issue.ID)
		proLeft := projectLen - len(issue.Project.Name)

		status := issue.Status.Name + strings.Repeat(" ", sLeft) // strconv.Itoa(sLeft)
		idPad := strings.Repeat(" ", iLeft)
		project := issue.Project.Name + strings.Repeat(" ", proLeft)

		fmt.Printf("%d%s  %s  %s  %s\n", issue.ID, idPad, status, project, issue.Subject)
	}
	fmt.Printf("--- Issues %d to %d (Total %d) ----\n",
		issues.Offset,
		issues.Limit,
		issues.TotalCount,
	)
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
