package print

import (
	"github.com/jedib0t/go-pretty/text"
)

type Color int

const (
	None Color = iota
	ID
	Blue
	Green
)

func (s Color) Color(value string) string {
	switch s {
	case ID:
		return text.FgGreen.Sprint(value)
	case Blue:
		return text.FgBlue.Sprint(value)
	case Green:
		return text.FgGreen.Sprint(value)
	}
	return value
}
