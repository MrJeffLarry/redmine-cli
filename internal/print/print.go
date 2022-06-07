package print

import (
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
)

func PrintDebug(r *config.Red_t, status int, body string) {
	if !r.Debug {
		return
	}
	fmt.Println("------------- DEBUG INFO - BEGIN ----------------")
	fmt.Println(body)
	fmt.Println("------------- DEBUG INFO - END ------------------")
}

func PrintRowHead(r *config.Red_t, format string, a ...any) {
	fmt.Printf(format, a...)
}

func PrintRow(r *config.Red_t, format string, a ...any) {
	fmt.Printf(format, a...)
}
