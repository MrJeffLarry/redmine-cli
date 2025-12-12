package cmd

import (
	"fmt"
	"os"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func cmdCompletion(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion -s <shell>",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for red-cli.

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
`,

		Run: func(cmd *cobra.Command, args []string) {
			shell, _ := cmd.Flags().GetString("shell")
			// flag takes precedence; fall back to positional arg for compatibility
			if shell == "" {
				if len(args) > 0 {
					shell = args[0]
				} else {
					// no shell provided; show help instead of plain error
					_ = cmd.Help()
					return
				}
			}

			switch shell {
			case "bash":
				_ = r.Cmd.GenBashCompletion(os.Stdout)
			case "zsh":
				_ = r.Cmd.GenZshCompletion(os.Stdout)
			case "fish":
				_ = r.Cmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = r.Cmd.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				// unsupported shell: show help to guide the user
				fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
				_ = cmd.Help()
			}
		},
	}
	cmd.Flags().StringP("shell", "s", "", "Shell type: {bash|zsh|fish|powershell}")
	return cmd
}
