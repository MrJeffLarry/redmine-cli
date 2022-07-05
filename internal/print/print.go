package print

import (
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/jedib0t/go-pretty/text"
)

func PrintDebug(r *config.Red_t, status int, body string) {
	if !r.Debug {
		return
	}
	fmt.Println(text.FgBlack.Sprintf("DEBUG: %s", body))
}

func Debug(r *config.Red_t, status int, body string) {
	if !r.Debug {
		return
	}
	fmt.Println(text.FgBlack.Sprintf("DEBUG: %s", body))
}

func PrintRowHead(r *config.Red_t, format string, a ...any) {
	fmt.Printf(format, a...)
}

func PrintRow(r *config.Red_t, format string, a ...any) {
	fmt.Printf(format, a...)
}

func PrintID(id int64) string {
	return text.FgGreen.Sprint("#", id)
}

func PrintTimeAgo(ago string) string {
	return ago
}

func Error(format string, a ...any) {
	fmt.Println(text.FgRed.Sprintf(format, a...))
}

func OK(format string, a ...any) {
	fmt.Println(text.FgGreen.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	fmt.Printf(format, a...)
}
