package issue

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"github.com/spf13/cobra"
)

const (
	FLAG_ORDER_DESC = "desc"

	FLAG_SORT   = "sort"
	FLAG_SORT_P = "s"

	FLAG_LIMIT   = "limit"
	FLAG_LIMIT_P = "l"

	FLAG_OFFSET   = "offset"
	FLAG_OFFSET_P = "o"

	FLAG_PAGE   = "page"
	FLAG_PAGE_P = "p"
)

func countDigi(i int64) (count int) {
	for i > 0 {
		i = i / 10
		count++
	}
	return
}

func parseFlags(cmd *cobra.Command, path string) string {
	FLAG_SORT_FIELDS := []string{"id", "status", "project", "subject"}

	limit, _ := cmd.Flags().GetInt(FLAG_LIMIT)
	offset, _ := cmd.Flags().GetInt(FLAG_OFFSET)
	page, _ := cmd.Flags().GetInt(FLAG_PAGE)
	sort, _ := cmd.Flags().GetString(FLAG_SORT)
	order, _ := cmd.Flags().GetBool(FLAG_ORDER_DESC)

	if page > 0 {
		path += "offset=" + strconv.Itoa(page*(limit)) + "&"
	} else {
		path += "offset=" + strconv.Itoa(offset) + "&"
	}
	path += "limit=" + strconv.Itoa(limit) + "&"

	if util.Contains(FLAG_SORT_FIELDS, sort) {
		path += "sort=" + sort
		if order {
			path += ":desc"
		}
		path += "&"
	}
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
		"ID "+strings.Repeat(" ", int(math.Abs(float64(idLen-len("ID"))))),
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

		fmt.Printf("%s%s  %s  %s  %s\n", print.PrintID(issue.ID), idPad, status, project, issue.Subject)
	}
	fmt.Printf("--- Issues %d to %d (Total %d) ----\n",
		issues.Offset,
		issues.Offset+issues.Limit,
		issues.TotalCount,
	)
}

func cmdIssueList(r *config.Red_t) *cobra.Command {
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
			displayListGET(r, cmd, "/issues.json?assigned_to_id=me&")
		},
	})

	cmd.PersistentFlags().Bool(FLAG_ORDER_DESC, false, "desc")
	cmd.PersistentFlags().StringP(FLAG_SORT, FLAG_SORT_P, "", "Sort field: ID, Status, Project, Subject")
	cmd.PersistentFlags().IntP(FLAG_PAGE, FLAG_PAGE_P, 0, "List 25 issues per page (uses limit and offset)")
	cmd.PersistentFlags().IntP(FLAG_LIMIT, FLAG_LIMIT_P, 25, "Limit number of issues per page")
	cmd.PersistentFlags().IntP(FLAG_OFFSET, FLAG_OFFSET_P, 0, "skip this number of issues")

	return cmd
}
