package issue

import (
	"encoding/json"
	"fmt"
	"net/url"
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
	FLAG_ORDER_ASC  = "asc"

	FLAG_SORT   = "sort"
	FLAG_SORT_P = "s"

	FLAG_SEARCH   = "search"
	FLAG_SEARCH_P = "q"

	FLAG_LIMIT   = "limit"
	FLAG_LIMIT_P = "l"

	FLAG_OFFSET   = "offset"
	FLAG_OFFSET_P = "o"

	FLAG_PAGE   = "page"
	FLAG_PAGE_P = "p"

	HEAD_ID       = "ID"
	HEAD_STATUS   = "STATUS"
	HEAD_PRIORITY = "PRIORITY"
	HEAD_PROJECT  = "PROJECT"
	HEAD_SUBJECT  = "SUBJECT"
)

func countDigi(i int64) (count int) {
	for i > 0 {
		i = i / 10
		count++
	}
	return
}

func parseFlags(r *config.Red_t, cmd *cobra.Command, path string) string {
	FLAG_SORT_FIELDS := []string{"id", "status", "priority", "project", "subject"}

	limit, _ := cmd.Flags().GetInt(FLAG_LIMIT)
	offset, _ := cmd.Flags().GetInt(FLAG_OFFSET)
	page, _ := cmd.Flags().GetInt(FLAG_PAGE)
	sort, _ := cmd.Flags().GetString(FLAG_SORT)
	order, _ := cmd.Flags().GetBool(FLAG_ORDER_ASC)
	search, _ := cmd.Flags().GetString(FLAG_SEARCH)

	if r.RedmineProjectID > 0 && !r.All {
		path += "project_id=" + strconv.Itoa(r.RedmineProjectID) + "&"
	}

	if len(search) > 0 {
		path += "subject=" + url.QueryEscape(search) + "&"
	}

	if page > 0 {
		path += "offset=" + strconv.Itoa(page*(limit)) + "&"
	} else {
		path += "offset=" + strconv.Itoa(offset) + "&"
	}
	path += "limit=" + strconv.Itoa(limit) + "&"

	if util.Contains(FLAG_SORT_FIELDS, sort) {
		path += "sort=" + sort
		if !order {
			path += ":desc"
		}
		path += "&"
	} else {
		path += "sort=priority"
		if !order {
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
	idLen := len(HEAD_ID)
	statusLen := len(HEAD_STATUS)
	priorityLen := len(HEAD_PRIORITY)
	projectLen := len(HEAD_PROJECT)
	subjectLen := len(HEAD_SUBJECT)

	issues := issues{}

	path = parseFlags(r, cmd, path)

	print.Debug(r, path)

	if body, status, err = api.ClientGET(r, path); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.Debug(r, "%d %s", status, string(body))

	if err := json.Unmarshal(body, &issues); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	for _, issue := range issues.Issues {
		iLen := countDigi(issue.ID)
		sLen := len(issue.Status.Name)
		pLen := len(issue.Priority.Name)
		proLen := len(issue.Project.Name)
		subLen := len(issue.Subject)

		if iLen > idLen {
			idLen = iLen
		}

		if sLen > statusLen {
			statusLen = sLen
		}

		if pLen > priorityLen {
			priorityLen = pLen
		}

		if proLen > projectLen {
			projectLen = proLen
		}

		if subLen > subjectLen {
			subjectLen = subLen
		}

	}

	fmt.Printf("%s   %s  %s  %s  %s\n",
		HEAD_ID+strings.Repeat(" ", idLen-len(HEAD_ID)),
		HEAD_STATUS+strings.Repeat(" ", statusLen-len(HEAD_STATUS)),
		HEAD_PRIORITY+strings.Repeat(" ", priorityLen-len(HEAD_PRIORITY)),
		HEAD_PROJECT+strings.Repeat(" ", projectLen-len(HEAD_PROJECT)),
		HEAD_SUBJECT,
	)

	for _, issue := range issues.Issues {
		iLeft := idLen - countDigi(issue.ID)
		sLeft := statusLen - len(issue.Status.Name)
		pLeft := priorityLen - len(issue.Priority.Name)
		proLeft := projectLen - len(issue.Project.Name)

		id := print.PrintID(issue.ID) + strings.Repeat(" ", iLeft)
		status := issue.Status.Name + strings.Repeat(" ", sLeft)
		priority := issue.Priority.Name + strings.Repeat(" ", pLeft)
		project := issue.Project.Name + strings.Repeat(" ", proLeft)

		fmt.Printf("%s  %s  %s  %s  %s\n",
			id,
			status,
			priority,
			project,
			issue.Subject,
		)
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
		Long:  "List all issues and ignores project ID",
		Run: func(cmd *cobra.Command, args []string) {
			r.RedmineProjectID = 0 // ignore ID if we want too see all
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

	cmd.PersistentFlags().Bool(FLAG_ORDER_ASC, false, "Ascend order")
	cmd.PersistentFlags().StringP(FLAG_SORT, FLAG_SORT_P, "", "Sort field: id, status, priority, project, subject")
	cmd.PersistentFlags().StringP(FLAG_SEARCH, FLAG_SEARCH_P, "", "Search in subject field")
	cmd.PersistentFlags().IntP(FLAG_PAGE, FLAG_PAGE_P, 0, "List 25 issues per page (uses limit and offset)")
	cmd.PersistentFlags().IntP(FLAG_LIMIT, FLAG_LIMIT_P, 25, "Limit number of issues per page")
	cmd.PersistentFlags().IntP(FLAG_OFFSET, FLAG_OFFSET_P, 0, "skip this number of issues")

	return cmd
}
