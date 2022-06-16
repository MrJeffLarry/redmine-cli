package main

import (
	"fmt"
	"os"
)

var Version = "0.0.0-dev"
var GitCommit = ""
var BuildTime = ""

func main() {
	if err := CmdInit(Version, GitCommit, BuildTime); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
