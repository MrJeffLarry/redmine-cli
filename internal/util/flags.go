package util

import (
	"strconv"

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
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(FLAG_ORDER_ASC, false, "Ascend order")
	cmd.PersistentFlags().StringP(FLAG_SORT, FLAG_SORT_P, "", "Sort field")
	//	cmd.PersistentFlags().StringP(FLAG_SEARCH, FLAG_SEARCH_P, "", "Search in subject field")
	cmd.PersistentFlags().IntP(FLAG_PAGE, FLAG_PAGE_P, 0, "List 25 objects per page (uses limit and offset)")
	cmd.PersistentFlags().IntP(FLAG_LIMIT, FLAG_LIMIT_P, 25, "Limit number of objects per page")
	cmd.PersistentFlags().IntP(FLAG_OFFSET, FLAG_OFFSET_P, 0, "skip this number of objects")
}

func ParseFlags(cmd *cobra.Command, projectID int, sortFields []string) string {
	path := ""

	limit, _ := cmd.Flags().GetInt(FLAG_LIMIT)
	offset, _ := cmd.Flags().GetInt(FLAG_OFFSET)
	sort, _ := cmd.Flags().GetString(FLAG_SORT)
	order, _ := cmd.Flags().GetBool(FLAG_ORDER_ASC)
	page, _ := cmd.Flags().GetInt(FLAG_PAGE)
	//	search, _ := cmd.Flags().GetString(FLAG_SEARCH)

	if projectID > 0 {
		path += "project_id=" + strconv.Itoa(projectID) + "&"
	}

	//	if len(search) > 0 {
	//		path += "subject=" + url.QueryEscape(search) + "&"
	//	}

	if page > 0 {
		path += "offset=" + strconv.Itoa(page*(limit)) + "&"
	} else {
		path += "offset=" + strconv.Itoa(offset) + "&"
	}
	path += "limit=" + strconv.Itoa(limit) + "&"

	if Contains(sortFields, sort) {
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
