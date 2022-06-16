package project

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func setProject(r *config.Red_t, cmd *cobra.Command) {
	id := cmd.Flags().Arg(0)

	r.SetProject(id)
	if err := r.Save(); err != nil {
		print.Error("%s [%s]", "Could not save project, please verify permissions", err)
		return
	}
	print.OK("Project globally set to %s", id)
}

func cmdProjectSet(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set project",
		Long:  "set all projects",
		Run: func(cmd *cobra.Command, args []string) {
			setProject(r, cmd)
		},
	}

	return cmd
}
