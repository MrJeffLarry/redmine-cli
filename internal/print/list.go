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
	Parent      int
	ParentSize  int
	ContentSize int
	Content     string
}

type List struct {
	maxLens     []int
	headers     []string
	rows        [][]Column
	Parent      int
	OldParent   int
	ParentLevel int
	Offset      int
	Limit       int
	TotalCount  int
}

func NewList(header ...string) *List {
	list := &List{}

	for _, field := range header {
		list.maxLens = append(list.maxLens, utf8.RuneCountInString(field))
		list.headers = append(list.headers, field)
	}
	return list
}

func (l *List) AddRow(row ...Column) {
	if len(row) != len(l.headers) {
		Error("Ohhh no! Darn! The number of columns does not match headers, Please report this to the developers")
		os.Exit(0)
	}
	for i, field := range row {
		if field.ParentPad {
			if field.Parent > 0 && field.Parent == l.Parent {
				// same level do nothing
			} else if field.Parent > 0 && field.Parent != l.Parent {
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
		row[i].ContentSize = utf8.RuneCountInString(field.Content) + (row[i].ParentSize * 2)
		if row[i].ContentSize > l.maxLens[i] {
			l.maxLens[i] = row[i].ContentSize
		}
	}
	l.rows = append(l.rows, row)
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

func (l *List) Render() {
	for i, head := range l.headers {
		fmt.Printf("%s  %s", head, strings.Repeat(" ", l.maxLens[i]-utf8.RuneCountInString(head)))
	}
	fmt.Printf("\n")
	for _, row := range l.rows {
		for i, field := range row {
			pad := l.maxLens[i]
			pad -= field.ContentSize
			fmt.Printf("%s%s  %s",
				strings.Repeat("‣ ", field.ParentSize),
				field.FgColor.Color(field.Content),
				strings.Repeat(" ", pad),
			)
		}
		fmt.Printf("\n")
	}
	if l.TotalCount == 0 {
		l.TotalCount = len(l.rows)
	}
	if l.Limit == 0 {
		l.Limit = l.TotalCount
	}
	if l.Offset == 0 {
		l.Offset = 0
	}

	fmt.Println(text.FgHiBlack.Sprintf("- - - - %d to %d (Total %d) - - - -",
		l.Offset,
		l.Offset+l.Limit,
		l.TotalCount,
	))
}
