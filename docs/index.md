# Documentation

[**Home**](../README.md) | Index

## Overview

`red-cli` is a small command-line client for Redmine. It provides subcommands to authenticate, list and manage projects and issues, and to configure local settings.

**NEW**: Red-cli now supports multi-instance management! Work with multiple Redmine servers simultaneously using the `--rid` flag. See the [Multi-Instance Guide](multi-instance.md) for details.

Full reference
----------

Read the full reference documentation for all commands and options below.

- [red-cli](red-cli.md) — Top-level CLI reference
- [config](config.md) — Configuration options and files
- [multi-instance](multi-instance.md) — Multi-instance support guide

Completion
----------

Use `red-cli completion -s <shell>` to generate completion scripts for bash, zsh, fish and PowerShell. See [completion](red-cli_completion.md) for details and install instructions.

Examples
--------

- Login using Username and password or API key:
	```bash
	red-cli auth login
	```

- Login to multiple instances:
	```bash
	red-cli auth login --rid prod
	red-cli auth login --rid staging
	```

- List projects:
	```bash
	red-cli project list
	```

- List projects from a specific instance:
	```bash
	red-cli project list --rid staging
	```

- Create an issue (opens editor for description):
	```bash
	red-cli issue create --project 42
	```

- Create an issue in a specific instance:
	```bash
	red-cli issue create --project 42 --rid prod
	```

- Add a note to an issue:
	```bash
	red-cli issue note 123 -m "This is a note"
	```


Want improvements?
-------------------

If you'd like dynamic shell completions (project/issue name completion), editor-based workflows, or example workflows for CI, open an issue or send a patch. This documentation is intentionally concise — see individual pages for full command reference.
