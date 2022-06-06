package main

import (
	"errors"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/cmd/issues"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/login"
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/users"
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

func CmdInit(Version string) error {

	r := config.InitConfig()

	cmd := &cobra.Command{
		Use:           "red <command> <subcommand> [flags]",
		Short:         "Redmine CLI",
		Long:          `Redmine CLI for integration with Redmine API`,
		Version:       Version,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run:           func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// require that the user is authenticated before running most commands
		if IsAuthCmd(cmd) && r.IsConfigBad() {
			fmt.Println("Redmine CLI (red) v" + Version)
			fmt.Println("")
			fmt.Println("You are not logged in, Please run `red login`")
			fmt.Println("")
			return errors.New("Not authenticated")
		}

		return nil
	}

	//	cmd.AddCommand(NewCmdVersion(Version))
	cmd.AddCommand(issues.NewCmdIssues(r))
	cmd.AddCommand(users.NewCmdUsers(r))
	cmd.AddCommand(login.NewCmdLogin(r))

	cmd.Execute()

	return nil
}
