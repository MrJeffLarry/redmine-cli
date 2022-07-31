package main

import (
	"fmt"
	"os"
)

var version = "0.0.0-dev"
var commit = ""
var date = ""

func main() {
	if err := CmdInit(version, commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
