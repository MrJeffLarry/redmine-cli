package issues

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

/*
issue - A hash of the issue attributes:
	- project_id
	- tracker_id - Bug=1, Feature=2, Support=3, Journalföring=4
	- status_id
	- priority_id
	- subject
	- description
	- category_id
	- fixed_version_id - ID of the Target Versions (previously called 'Fixed Version' and still referred to as such in the API)
	- assigned_to_id - ID of the user to assign the issue to (currently no mechanism to assign by name)
	- parent_issue_id - ID of the parent issue
	- custom_fields - See Custom fields
	- watcher_user_ids - Array of user ids to add as watchers (since 2.3.0)
	- is_private - Use true or false to indicate whether the issue is private or not
	- estimated_hours - Number of hours estimated for issue
*/

func NewCmdIssues(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issues",
		Short:   "issues",
		Long:    "issues",
		Aliases: []string{"i"},
		Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
	}

	cmd.AddCommand(cmdIssuesList(r))

	return cmd
}
