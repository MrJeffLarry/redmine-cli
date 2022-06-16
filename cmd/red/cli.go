package main

import (
	"errors"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/cmd/auth"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/issue"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/project"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/user"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
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

func CmdInit(Version, GitCommit, BuildTime string) error {

	r := config.InitConfig()

	cmd := &cobra.Command{
		Use:           "red <command> <subcommand> [flags]",
		Short:         "Redmine CLI",
		Long:          `Redmine CLI for integration with Redmine API`,
		Version:       Version + "\nGit Commit: " + GitCommit + "\nBuild time: " + BuildTime,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run:           func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Debug flag check for more info
		r.Debug, _ = cmd.Flags().GetBool(config.DEBUG_FLAG)

		// require that the user is authenticated before running most commands
		if IsAuthCmd(cmd) && r.IsConfigBad() {
			fmt.Println("Redmine CLI (red) v" + Version)
			fmt.Println("")
			fmt.Println("You are not logged in, Please run `red auth login`")
			fmt.Println("")
			return errors.New("Not authenticated")
		}

		return nil
	}

	cmd.PersistentFlags().BoolP(config.DEBUG_FLAG, config.DEBUG_FLAG_S, false, "Show debug info and raw response")

	cmd.AddCommand(issue.NewCmdIssue(r))
	cmd.AddCommand(project.NewCmdProject(r))
	cmd.AddCommand(user.NewCmdUser(r))
	cmd.AddCommand(auth.NewCmdAuth(r))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}

	return nil
}
