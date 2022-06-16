package project

import (
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func setProject(r *config.Red_t, cmd *cobra.Command) {
	var id int
	var err error

	if id, err = strconv.Atoi(cmd.Flags().Arg(0)); err != nil {
		print.Error("ID is not an valid number, please use `project set 1` for project id 1")
		return
	}

	r.SetProjectID(id)
	if err := r.Save(); err != nil {
		print.Error("%s [%s]", "Could not save project, please verify permissions", err)
		return
	}
	print.OK("Project globally set to ID #%d", id)
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
