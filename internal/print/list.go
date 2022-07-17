package print

import (
	"fmt"
	"strings"
)

type List struct {
	maxLens []int
	headers []string
	rows    [][]string
}

func NewList(header ...string) *List {
	list := &List{}

	for _, field := range header {
		list.maxLens = append(list.maxLens, len(field))
		list.headers = append(list.headers, field)
	}
	return list
}

func (l *List) AddRow(row ...string) {
	for i, field := range row {
		if len(field) > l.maxLens[i] {
			l.maxLens[i] = len(field)
		}
	}
	l.rows = append(l.rows, row)
}

func (l *List) Render() {
	for i, head := range l.headers {
		fmt.Printf("%s %s:%d-%d", head, strings.Repeat(" ", l.maxLens[i]-len(head)), len(head), l.maxLens[i])
	}
	fmt.Printf("\n")
	for _, row := range l.rows {
		for i, field := range row {
			fmt.Printf("%s %s", field, strings.Repeat(" ", l.maxLens[i]-len(field)))
		}
		fmt.Printf("\n")
	}
}

/*
func TableBody(format, headersLen []int, )

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
}*/
