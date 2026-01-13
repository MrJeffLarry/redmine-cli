# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli completion

## red-cli completion

Generate shell completion script

### Synopsis

Generate shell completion script for red-cli.

To load completions:

Bash:

  First, ensure that you install bash-completion using your package manager.

  After, add this to your ~/.bash_profile:

  $ eval "$(red-cli completion -s bash)"

Zsh:

  $ red-cli completion -s zsh > /usr/local/share/zsh/site-functions/_red_cli

  Ensure that the following is present in your ~/.zshrc:

  autoload -U compinit
  compinit -i

Fish:

  Generate a red-cli.fish completion script:

  $ red-cli completion -s fish > ~/.config/fish/completions/red-cli.fish

PowerShell:

  Open your profile script with:

  $ mkdir -Path (Split-Path -Parent $profile) -ErrorAction SilentlyContinue
  $ notepad $profile

  Add the line and save the file:

  $ Invoke-Expression -Command $(red-cli completion -s powershell | Out-String)


```
red-cli completion -s <shell> [flags]
```

### Options

```
  -h, --help           help for completion
  -s, --shell string   Shell type: {bash|zsh|fish|powershell}
```

### Options inherited from parent commands

```
      --all          Ignore project-id
  -d, --debug        Show debug info and raw response
      --rid string   Redmine instance ID (for multi-instance support)
```

### SEE ALSO

* [red-cli](./red-cli.md)	 - Redmine CLI

