package issue

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdIssue(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "issue",
		Long:    "issue",
		Aliases: []string{"i"},
		Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdIssueList(r))
	cmd.AddCommand(cmdIssueView(r))

	return cmd
}
