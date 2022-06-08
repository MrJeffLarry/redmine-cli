package issues

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdIssues(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "issue",
		Long:    "issue",
		Aliases: []string{"i"},
		Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdIssuesList(r))
	cmd.AddCommand(cmdIssuesView(r))

	return cmd
}
