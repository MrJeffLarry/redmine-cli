package main

import (
	"fmt"
	"os"

	"github.com/MrJeffLarry/redmine-cli/internal/cmd"
)

var version = "0.0.0-dev"
var commit = ""
var date = ""

func main() {
	r := cmd.CmdInit(version, commit, date)
	if err := r.Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
