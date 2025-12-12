package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/cmd"
	"github.com/spf13/cobra/doc"
)

var version = "0.0.0-dev"
var commit = ""
var date = ""

func main() {
	r := cmd.CmdInit(version, commit, date)

	r.Cmd.DisableAutoGenTag = true

	const fmTemplate = `# Documentation
[**Home**](../README.md) | [**Index**](index.md) | %s

`

	filePrepender := func(filename string) string {
		base := filename[strings.LastIndex(filename, "/")+1:]
		base = strings.TrimSuffix(base, ".md")
		base = strings.ReplaceAll(base, "_", " ")
		return fmt.Sprintf(fmTemplate, base)
	}

	linkHandler := func(name string) string {
		return "./" + name
	}

	err := doc.GenMarkdownTreeCustom(r.Cmd, "./docs", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
