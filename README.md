# Redmine-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/MrJeffLarry/redmine-cli)](https://goreportcard.com/report/github.com/MrJeffLarry/redmine-cli)

Redmine-cli (Command Line Interface) is a software program that allows you to interact with the Redmine project management system using the command line.
With a Redmine CLI tool, you can perform a wide range of tasks, such as creating and managing issues, list projects, and more, all from the comfort of your terminal or command prompt.

Supports redmine versions

* 5.x
* 4.x

## Install

### MacOS

**brew**

```bash
brew tap mrjefflarry/redmine-cli https://github.com/mrjefflarry/redmine-cli
brew install mrjefflarry/redmine-cli/red-cli
```

### Windows

**scoop**

```powershell
scoop bucket add org https://github.com/mrjefflarry/redmine-cli.git
scoop install red-cli
```

### Linux

**apt**

```bash
curl -s --compressed "https://redmine-cli.hagerman.io/apt/public_key.gpg" | sudo apt-key add -
sudo curl -s --compressed -o /etc/apt/sources.list.d/redmine-cli.list "https://redmine-cli.hagerman.io/apt/redmine-cli.list"
sudo apt update
sudo apt install red-cli
```

**ArchLinux**

Install from the [redmine-cli-go AUR package](https://aur.archlinux.org/packages/redmine-cli-go)

```bash
yay redmine-cli-go
```

## Usage

```
> red-cli -h
Redmine CLI for integration with Redmine API

Usage:
  red-cli <command> <subcommand> [flags]
  red-cli [command]

Available Commands:
  auth        auth to Redmine
  completion  Generate the autocompletion script for the specified shell
  config      get or set config
  help        Help about any command
  issue       issue
  project     project
  user        Get users info

Flags:
      --all       Ignore project-id
  -d, --debug     Show debug info and raw response
  -h, --help      help for red
  -v, --version   version for red

Use "red-cli [command] --help" for more information about a command.
```

Documentation can be found [here](./docs/config.md)
