package project

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func setProject(r *config.Red_t, cmd *cobra.Command) {
	id := cmd.Flags().Arg(0)

	if len(id) == 0 {
		print.Error("Please specify what project identifier you would like to use globally, usage: set [id]")
		return
	}
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
