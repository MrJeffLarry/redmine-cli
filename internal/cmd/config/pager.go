package config

import (
	"errors"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func getPager(r *config.Red_t, cmd *cobra.Command) (string, error) {
	view := cmd.Flags().Arg(0)

	if len(view) <= 0 {
		print.Error("Missing pager exec, please specify what program to use to view")
		return "", errors.New("")
	}
	return view, nil
}

func setGlobalPager(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set pager",
		Long:  "this will set pager globally in ~/.red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id string
			var err error
			if id, err = getPager(r, cmd); err != nil {
				return
			}

			r.SetPager(id)
			if err := r.Save(); err != nil {
				print.Error("%s [%s]", "Could not save pager, please verify permissions", err)
				return
			}
			print.OK("Pager set to [%s]", id)
		},
	}
	return cmd
}

func globalPager(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pager",
		Short: "set or get pager",
		Long:  "set or get what pager to use globally in ~/.red/config.json",
	}
	cmd.AddCommand(setGlobalPager(r))
	return cmd
}
