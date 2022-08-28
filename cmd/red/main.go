package main

import (
	"fmt"
	"os"

	"github.com/MrJeffLarry/redmine-cli/internal/editor"
)

var version = "0.0.0-dev"
var commit = ""
var date = ""

func main() {
	editor.StartEdit("")
	if err := CmdInit(version, commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
