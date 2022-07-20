package print

import (
	"fmt"
	"strings"
)

type Colum struct {
	FgColor     Color
	BgColor     Color
	ParentPad   bool
	Parent      string
	ParentSize  int
	ContentSize int
	Content     string
}

type List struct {
	maxLens     []int
	headers     []string
	rows        [][]Colum
	Parent      string
	OldParent   string
	ParentLevel int
}

func NewList(header ...string) *List {
	list := &List{}

	for _, field := range header {
		list.maxLens = append(list.maxLens, len(field))
		list.headers = append(list.headers, field)
	}
	return list
}

func (l *List) AddRow(row ...Colum) {
	for i, field := range row {
		if field.ParentPad {
			if len(field.Parent) > 0 && field.Parent == l.Parent {
				// same level do nothing
			} else if len(field.Parent) > 0 && field.Parent != l.Parent {
				if l.OldParent == field.Parent {
					l.ParentLevel--
				} else {
					l.ParentLevel++
				}
				l.OldParent = l.Parent
				l.Parent = field.Parent
			} else {
				l.Parent = field.Parent
				l.ParentLevel = 0
			}
			row[i].ParentSize = l.ParentLevel
		}
		row[i].ContentSize = len(field.Content) + (row[i].ParentSize * 2)
		if row[i].ContentSize > l.maxLens[i] {
			l.maxLens[i] = row[i].ContentSize
		}
	}
	l.rows = append(l.rows, row)
}

func (l *List) Render() {
	for i, head := range l.headers {
		fmt.Printf("%s %s", head, strings.Repeat(" ", l.maxLens[i]-len(head)))
	}
	fmt.Printf("\n")
	for _, row := range l.rows {
		for i, field := range row {
			pad := l.maxLens[i]
			pad -= field.ContentSize
			fmt.Printf("%s%s %s",
				strings.Repeat("‣ ", field.ParentSize),
				field.FgColor.Color(field.Content),
				strings.Repeat(" ", pad))
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
