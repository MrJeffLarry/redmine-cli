package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/cmd/auth"
	cmdConfig "github.com/MrJeffLarry/redmine-cli/internal/cmd/config"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/issue"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/project"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/user"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/spf13/cobra"
)

func IsAuthCmd(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help":
		return false
	case "version":
		return false
	case "login":
		return false
	}
	return true
}

func CmdInit(Version, GitCommit, BuildTime string) *config.Red_t {

	r := config.InitConfig()

	r.Term = terminal.New(nil, nil, nil)
	version := Version
	if strings.TrimSpace(GitCommit) != "" {
		version += "\nGit Commit: " + GitCommit
	}

	if strings.TrimSpace(BuildTime) != "" {
		version += "\nBuild Time: " + BuildTime
	}

	r.Cmd = &cobra.Command{
		Use:           "red-cli <command> <subcommand> [flags]",
		Short:         "Redmine CLI",
		Long:          `Redmine CLI for integration with Redmine API`,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run:           func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	r.Cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Debug flag check for more info
		r.Debug, _ = cmd.Flags().GetBool(config.DEBUG_FLAG)
		r.All, _ = cmd.Flags().GetBool(config.ALL_FLAG)

		// require that the user is authenticated before running most commands
		if IsAuthCmd(cmd) && r.IsConfigBad() {
			fmt.Println("Redmine CLI (red-cli) v" + Version)
			fmt.Println("")
			fmt.Println("You are not logged in, Please run `red-cli auth login`")
			fmt.Println("")
			return errors.New("Not authenticated")
		}

		return nil
	}

	r.Cmd.PersistentFlags().BoolP(config.DEBUG_FLAG, config.DEBUG_FLAG_S, false, "Show debug info and raw response")
	r.Cmd.PersistentFlags().Bool(config.ALL_FLAG, false, "Ignore project-id")

	r.Cmd.AddCommand(issue.NewCmdIssue(r))
	r.Cmd.AddCommand(project.NewCmdProject(r))
	r.Cmd.AddCommand(user.NewCmdUser(r))
	r.Cmd.AddCommand(auth.NewCmdAuth(r))
	r.Cmd.AddCommand(cmdConfig.NewCmdConfig(r))
	r.Cmd.AddCommand(cmdCompletion(r))

	return r
}
