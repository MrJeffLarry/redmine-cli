package config

import (
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func getProjectID(r *config.Red_t, cmd *cobra.Command) (int, error) {
	var id int
	var err error

	if id, err = strconv.Atoi(cmd.Flags().Arg(0)); err != nil {
		print.Error("ID is not an valid number, please use `project set 1` for project id 1")
		return 0, err
	}
	return id, nil
}

func setLocalProject(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set local project",
		Long:  "set a local project under .red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			if id, err = getProjectID(r, cmd); err != nil {
				return
			}

			if err := r.SaveLocalProject(id); err != nil {
				print.Error("%s [%s]", "Could not save project, please verify permissions", err)
				return
			}
			print.OK("Project locally set to ID #%d", id)
		},
	}
	return cmd
}

func setGlobalProject(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set global project",
		Long:  "set a global project under ~/.red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			if id, err = getProjectID(r, cmd); err != nil {
				return
			}

			r.SetProjectID(id)
			if err := r.Save(); err != nil {
				print.Error("%s [%s]", "Could not save project, please verify permissions", err)
				return
			}
			print.OK("Project globally set to ID #%d", id)
		},
	}
	return cmd
}

func localProject(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "set or get local project",
		Long:  "set or get a local project under ~/.red/config.json",
	}
	cmd.AddCommand(setLocalProject(r))
	return cmd
}

func globalProject(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "set or get global project",
		Long:  "set or get a global project under ~/.red/config.json",
	}
	cmd.AddCommand(setGlobalProject(r))
	return cmd
}
