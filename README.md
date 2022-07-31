# Redmine-cli

Redmine CLI

## Install

### MacOS

**brew**

```bash
brew tap mrjefflarry/redmine-cli https://github.com/mrjefflarry/redmine-cli
brew install redmine-cli
```

### Windows

**scoop**

```powershell
scoop bucket add org https://github.com/mrjefflarry/redmine-cli.git
scoop install mrjefflarry/redmine-cli
```

### Linux

**snap**

```bash
snap install red-cli
```

## Usage

```
> red -h
Redmine CLI for integration with Redmine API

Usage:
  red <command> <subcommand> [flags]
  red [command]

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

Use "red [command] --help" for more information about a command.
```

### Usage for override global config

You can override the global config with a local config, for example if you use one repo for one project and other for other project you can create a folder inside current working directory **.red** inside that create a file called **config.json** this can then contain and override one or more options below

```bash
.red/config.json
```

contains 

```json
{
    "project-id": 12
}
```

this will then override the project

### Config

**Complete config list options**

```json
{
    "server": "https://redmine.example.com",
    "api-key": "randomkeyfromredmine",
    "user-id": 1,
    "project-id": 23,
}
```