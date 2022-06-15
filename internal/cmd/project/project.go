package project

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdProject(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Short:   "project",
		Long:    "Project",
		Aliases: []string{"p"},
		Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdProjectList(r))
	cmd.AddCommand(cmdProjectSet(r))

	return cmd
}
