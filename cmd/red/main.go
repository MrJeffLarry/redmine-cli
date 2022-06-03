package main

import (
	"fmt"
	"os"
)

func main() {
	appVersion := "0.0.1"

	if err := CmdInit(appVersion); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
