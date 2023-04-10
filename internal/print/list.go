package print

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/jedib0t/go-pretty/text"
)

type Column struct {
	FgColor     Color
	BgColor     Color
	ParentPad   bool
	ParentSize  int
	ContentSize int
	Content     string
}

type List struct {
	maxLens             []int
	headers             []string
	rows                []Row
	newRows             []Row
	Parent              int
	OldParent           int
	ParentLevel         int
	Offset              int
	Limit               int
	TotalCount          int
	ParentIssueGrouping bool
}

type Row struct {
	ID        int
	ParentID  int
	IgnorePad bool
	Solved    bool
	Columns   []Column
}

func NewList(header ...string) *List {
	list := &List{}

	list.ParentIssueGrouping = true

	for _, field := range header {
		list.maxLens = append(list.maxLens, utf8.RuneCountInString(field))
		list.headers = append(list.headers, field)
	}
	return list
}

func (l *List) SetParentIssueGrouping(v bool) {
	l.ParentIssueGrouping = v
}

func (l *List) AddRow(id int, parentID int, row ...Column) {
	if len(row) != len(l.headers) {
		Error("Ohhh no! Darn! The number of columns does not match headers, Please report this to the developers")
		os.Exit(0)
	}

	r := Row{}
	r.ID = id
	r.ParentID = parentID
	r.Columns = append(r.Columns, row...)
	l.rows = append(l.rows, r)
}

func (l *List) SetOffset(offset int) {
	l.Offset = offset
}

func (l *List) SetLimit(limit int) {
	l.Limit = limit
}

func (l *List) SetTotal(total int) {
	l.TotalCount = total
}

func child(l *List, parentI int) {
	for i, row := range l.rows {
		if l.rows[parentI].ID == row.ParentID {
			l.rows[i].Solved = true
			l.newRows = append(l.newRows, row)
			child(l, i)
		}
	}
}

func (l *List) Render() string {

	if l.ParentIssueGrouping {
		for i, row := range l.rows {
			if row.ParentID == 0 {
				l.rows[i].Solved = true
				l.newRows = append(l.newRows, row)
				child(l, i)
			}
		}

		for _, row := range l.rows {
			if !row.Solved {
				row.ParentID = 0
				l.newRows = append(l.newRows, row)
			}
		}
	}

	for _, row := range l.newRows {
		for i1, field := range row.Columns {
			if field.ParentPad && l.ParentIssueGrouping && !row.IgnorePad {
				if row.ParentID > 0 && row.ParentID == l.Parent {
					// same level do nothing
				} else if row.ParentID > 0 && row.ParentID != l.Parent {
					if l.OldParent == row.ParentID {
						l.ParentLevel--
					} else {
						l.ParentLevel++
					}
					l.OldParent = l.Parent
					l.Parent = row.ParentID
				} else {
					l.Parent = row.ParentID
					l.ParentLevel = 0
				}
				row.Columns[i1].ParentSize = l.ParentLevel
			}

			row.Columns[i1].ContentSize = utf8.RuneCountInString(field.Content) + (row.Columns[i1].ParentSize * 2)

			if row.Columns[i1].ContentSize > l.maxLens[i1] {
				l.maxLens[i1] = row.Columns[i1].ContentSize
			}
		}
	}

	output := ""

	for i, head := range l.headers {

		output += fmt.Sprintf("%s  %s", head, strings.Repeat(" ", l.maxLens[i]-utf8.RuneCountInString(head)))
	}

	output += "\n"

	for _, row := range l.newRows {
		for i, field := range row.Columns {
			parentSize := 0
			pad := l.maxLens[i]
			pad -= field.ContentSize

			if field.ParentSize > 0 {
				parentSize = field.ParentSize
			}

			output += fmt.Sprintf("%s%s  %s",
				strings.Repeat("‣ ", parentSize),
				field.FgColor.Color(field.Content),
				strings.Repeat(" ", pad),
			)
		}
		output += "\n"
	}
	if l.TotalCount == 0 {
		l.TotalCount = len(l.newRows)
	}
	if l.Limit == 0 {
		l.Limit = l.TotalCount
	}
	if l.Offset == 0 {
		l.Offset = 0
	}

	output += fmt.Sprintln(text.FgHiBlack.Sprintf("- - - - %d to %d (Total %d) - - - -",
		l.Offset,
		l.Offset+l.Limit,
		l.TotalCount,
	))

	return output
}
